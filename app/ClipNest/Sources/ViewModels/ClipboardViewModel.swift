import Foundation
import SwiftUI

@MainActor
final class ClipboardViewModel: ObservableObject {
    @Published var clips: [Clip] = []
    @Published var searchQuery = ""
    @Published var isConnected = false
    @Published var hoveredClip: Clip? = nil

    private let client = DaemonClient()
    private var broadcastTask: Task<Void, Never>?
    private var searchTask: Task<Void, Never>?
    func setHoveredClip(_ clip: Clip?) {
        guard let clip else { return }
        hoveredClip = clip
    }

    init() {
        Task { await start() }
    }

    private func start() async {
        await client.connect()

        // Poll connection status and load initial clips
        broadcastTask = Task {
            // Wait briefly for connection
            try? await Task.sleep(for: .milliseconds(500))
            let connected = await client.connected
            self.isConnected = connected
            if connected {
                await loadClips()
            }

            // Listen for broadcasts
            for await clip in client.broadcasts {
                self.isConnected = true
                // Insert at top, removing duplicate if exists
                self.clips.removeAll { $0.id == clip.id }
                self.clips.insert(clip, at: 0)
            }
            self.isConnected = false
        }

        // Periodically check connection status
        Task {
            while !Task.isCancelled {
                try? await Task.sleep(for: .seconds(2))
                let connected = await client.connected
                if self.isConnected != connected {
                    self.isConnected = connected
                    if connected {
                        await loadClips()
                    }
                }
            }
        }
    }

    func loadClips() async {
        do {
            let result = try await client.listClips()
            self.clips = result
        } catch {
            // Silently handle — UI shows disconnected state
        }
    }

    func search() {
        searchTask?.cancel()
        searchTask = Task {
            try? await Task.sleep(for: .milliseconds(300))
            guard !Task.isCancelled else { return }

            do {
                let query = searchQuery.trimmingCharacters(in: .whitespacesAndNewlines)
                if query.isEmpty {
                    let result = try await client.listClips()
                    self.clips = result
                } else {
                    let result = try await client.searchClips(query: query)
                    self.clips = result
                }
            } catch {
                // Silently handle
            }
        }
    }

    func copyClip(_ clip: Clip) {
        Task {
            do {
                try await client.copyClip(id: clip.id)
            } catch {
                // Silently handle
            }
        }
    }

    func togglePin(_ clip: Clip) {
        Task {
            do {
                if clip.pinned {
                    try await client.unpinClip(id: clip.id)
                } else {
                    try await client.pinClip(id: clip.id)
                }
                // Update local state
                if let index = clips.firstIndex(where: { $0.id == clip.id }) {
                    clips[index].pinned.toggle()
                }
            } catch {
                // Silently handle
            }
        }
    }

    func clearAll() {
        Task {
            do {
                try await client.clearClips()
                hoveredClip = nil
                clips.removeAll()
            } catch {
                // Silently handle
            }
        }
    }
}
