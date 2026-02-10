import Foundation

/// Outer wire envelope: {"type":"...", "data": ...}
struct RawSocketMessage: Sendable {
    let type: String
    let rawData: Data?

    init?(from jsonData: Data) {
        guard let obj = try? JSONSerialization.jsonObject(with: jsonData) as? [String: Any],
              let type = obj["type"] as? String else {
            return nil
        }
        self.type = type
        if let data = obj["data"] {
            self.rawData = try? JSONSerialization.data(withJSONObject: data)
        } else {
            self.rawData = nil
        }
    }
}

/// Response payload: {"success":true, "data":{"clips":[...], "count":N}, "error":"..."}
struct ResponseMessage: Codable, Sendable {
    let success: Bool
    let data: ClipListData?
    let error: String?
}

/// Clip list payload inside a response
struct ClipListData: Codable, Sendable {
    let clips: [Clip]
    let count: Int
}

/// Helpers for building outgoing command JSON
enum CommandBuilder {
    static func makeCommand(type: String, data: [String: Any]? = nil) -> Data? {
        var dict: [String: Any] = ["type": type]
        if let data {
            dict["data"] = data
        }
        guard let json = try? JSONSerialization.data(withJSONObject: dict) else {
            return nil
        }
        return json + [0x0A] // append newline
    }
}
