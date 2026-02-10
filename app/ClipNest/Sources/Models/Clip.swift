import Foundation

struct Clip: Identifiable, Codable, Sendable, Equatable {
    let id: Int64
    let content: String
    let type: String
    let timestamp: Int64
    var pinned: Bool

    var date: Date {
        Date(timeIntervalSince1970: TimeInterval(timestamp))
    }

    var preview: String {
        let trimmed = content.trimmingCharacters(in: .whitespacesAndNewlines)
        if trimmed.count > 120 {
            return String(trimmed.prefix(120)) + "..."
        }
        return trimmed
    }
}
