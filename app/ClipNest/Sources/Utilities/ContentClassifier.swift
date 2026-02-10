import Foundation

enum ContentClassifier {
    static func looksLikeCode(_ content: String) -> Bool {
        let lines = content.components(separatedBy: .newlines)
        var score = 0

        // Multiline + indentation: 3+ lines with at least one indented
        if lines.count >= 3 && lines.contains(where: { $0.hasPrefix(" ") || $0.hasPrefix("\t") }) {
            score += 1
        }

        // Brace pairs: 2+ { or } characters
        let braceCount = content.filter { $0 == "{" || $0 == "}" }.count
        if braceCount >= 2 {
            score += 1
        }

        // Trailing semicolons: 2+ lines ending in ;
        let semicolonLines = lines.filter { $0.trimmingCharacters(in: .whitespaces).hasSuffix(";") }.count
        if semicolonLines >= 2 {
            score += 1
        }

        // Programming keywords
        let keywords = ["func ", "def ", "class ", "import ", "var ", "let ", "const ",
                         "return ", "struct ", "enum ", "fn ", "pub ", "async ",
                         "interface ", "package "]
        let keywordMatches = keywords.filter { content.contains($0) }.count
        if keywordMatches >= 2 {
            score += 1
        }

        // Arrow operators
        if content.contains("=>") || content.contains("->") {
            score += 1
        }

        // Comment patterns
        let commentLines = lines.filter {
            let trimmed = $0.trimmingCharacters(in: .whitespaces)
            return trimmed.hasPrefix("//") || trimmed.hasPrefix("#") ||
                   trimmed.hasPrefix("/*") || trimmed.hasPrefix("*")
        }.count
        if commentLines >= 1 {
            score += 1
        }

        return score >= 2
    }

    static func classify(_ content: String) -> String {
        if looksLikeCode(content) {
            return "Code"
        } else if content.trimmingCharacters(in: .whitespacesAndNewlines).hasPrefix("http://") ||
                  content.trimmingCharacters(in: .whitespacesAndNewlines).hasPrefix("https://") {
            return "URL"
        } else {
            return "Text"
        }
    }
}
