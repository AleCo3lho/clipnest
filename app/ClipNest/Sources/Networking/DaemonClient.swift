import Foundation
import Network

actor DaemonClient {
    private let socketPath: String
    private var connection: NWConnection?
    private var buffer = Data()
    private var pendingContinuation: CheckedContinuation<ResponseMessage, Error>?
    private var broadcastContinuation: AsyncStream<Clip>.Continuation?
    private var isStarted = false

    private(set) var connected = false

    nonisolated let broadcasts: AsyncStream<Clip>

    init(socketPath: String = "/tmp/clipnest.sock") {
        self.socketPath = socketPath
        var continuation: AsyncStream<Clip>.Continuation?
        self.broadcasts = AsyncStream { continuation = $0 }
        self.broadcastContinuation = continuation
    }

    func connect() {
        guard !isStarted else { return }
        isStarted = true
        startConnection()
    }

    func disconnect() {
        isStarted = false
        connection?.cancel()
        connection = nil
        connected = false
        broadcastContinuation?.finish()
    }

    // MARK: - Connection

    private func startConnection() {
        let endpoint = NWEndpoint.unix(path: socketPath)
        let conn = NWConnection(to: endpoint, using: .tcp)
        self.connection = conn

        conn.stateUpdateHandler = { [weak self] state in
            guard let self else { return }
            Task { await self.handleStateChange(state) }
        }

        conn.start(queue: .global(qos: .userInitiated))
    }

    private func handleStateChange(_ state: NWConnection.State) {
        switch state {
        case .ready:
            connected = true
            buffer = Data()
            receiveLoop()

        case .failed, .cancelled:
            connected = false
            failPending(with: DaemonError.disconnected)
            if isStarted {
                scheduleReconnect()
            }

        case .waiting:
            connected = false

        default:
            break
        }
    }

    private func scheduleReconnect() {
        Task {
            try? await Task.sleep(for: .seconds(3))
            guard isStarted else { return }
            startConnection()
        }
    }

    // MARK: - Receiving

    private func receiveLoop() {
        guard let connection else { return }
        connection.receive(minimumIncompleteLength: 1, maximumLength: 65536) { [weak self] content, _, isComplete, error in
            guard let self else { return }
            Task {
                if let content {
                    await self.handleReceivedData(content)
                }
                if isComplete {
                    await self.handleStateChange(.cancelled)
                } else if error != nil {
                    await self.handleStateChange(.failed(error!))
                } else {
                    await self.receiveLoop()
                }
            }
        }
    }

    private func handleReceivedData(_ data: Data) {
        buffer.append(data)
        processBuffer()
    }

    private func processBuffer() {
        while let newlineIndex = buffer.firstIndex(of: 0x0A) {
            let lineData = buffer[buffer.startIndex..<newlineIndex]
            buffer = Data(buffer[(newlineIndex + 1)...])

            guard !lineData.isEmpty,
                  let message = RawSocketMessage(from: Data(lineData)) else {
                continue
            }

            dispatch(message)
        }
    }

    // MARK: - Dispatching

    private func dispatch(_ message: RawSocketMessage) {
        switch message.type {
        case "response":
            guard let rawData = message.rawData else {
                failPending(with: DaemonError.invalidResponse)
                return
            }
            do {
                let response = try JSONDecoder().decode(ResponseMessage.self, from: rawData)
                pendingContinuation?.resume(returning: response)
                pendingContinuation = nil
            } catch {
                failPending(with: error)
            }

        case "new_clip":
            guard let rawData = message.rawData else { return }
            if let clip = try? JSONDecoder().decode(Clip.self, from: rawData) {
                broadcastContinuation?.yield(clip)
            }

        default:
            break
        }
    }

    private func failPending(with error: Error) {
        pendingContinuation?.resume(throwing: error)
        pendingContinuation = nil
    }

    // MARK: - Sending

    private func sendCommand(type: String, data: [String: Any]? = nil) async throws -> ResponseMessage {
        guard let connection, connected else {
            throw DaemonError.notConnected
        }

        guard let payload = CommandBuilder.makeCommand(type: type, data: data) else {
            throw DaemonError.encodingFailed
        }

        return try await withCheckedThrowingContinuation { continuation in
            self.pendingContinuation = continuation
            connection.send(content: payload, completion: .contentProcessed { error in
                if let error {
                    Task { await self.failPending(with: error) }
                }
            })
        }
    }

    // MARK: - Convenience

    func listClips(limit: Int = 50) async throws -> [Clip] {
        let response = try await sendCommand(type: "list", data: ["limit": limit])
        return response.data?.clips ?? []
    }

    func searchClips(query: String, limit: Int = 50) async throws -> [Clip] {
        let response = try await sendCommand(type: "search", data: ["query": query, "limit": limit])
        return response.data?.clips ?? []
    }

    func copyClip(id: Int64) async throws {
        let response = try await sendCommand(type: "copy_clip", data: ["id": id])
        if !response.success {
            throw DaemonError.serverError(response.error ?? "unknown error")
        }
    }

    func pinClip(id: Int64) async throws {
        let response = try await sendCommand(type: "pin", data: ["id": id])
        if !response.success {
            throw DaemonError.serverError(response.error ?? "unknown error")
        }
    }

    func unpinClip(id: Int64) async throws {
        let response = try await sendCommand(type: "unpin", data: ["id": id])
        if !response.success {
            throw DaemonError.serverError(response.error ?? "unknown error")
        }
    }

    func getPinnedClips() async throws -> [Clip] {
        let response = try await sendCommand(type: "pins")
        return response.data?.clips ?? []
    }

    func clearClips() async throws {
        let response = try await sendCommand(type: "clear")
        if !response.success {
            throw DaemonError.serverError(response.error ?? "unknown error")
        }
    }
}

enum DaemonError: Error, LocalizedError {
    case notConnected
    case disconnected
    case encodingFailed
    case invalidResponse
    case serverError(String)

    var errorDescription: String? {
        switch self {
        case .notConnected: "Not connected to daemon"
        case .disconnected: "Disconnected from daemon"
        case .encodingFailed: "Failed to encode command"
        case .invalidResponse: "Invalid response from daemon"
        case .serverError(let msg): msg
        }
    }
}
