import SwiftUI

struct PreviewPanelView: View {
    let clip: Clip
    let onCopy: () -> Void

    private let maxDisplayLength = 5000

    private var contentType: String {
        ContentClassifier.classify(clip.content)
    }

    private var displayContent: String {
        if clip.content.count > maxDisplayLength {
            return String(clip.content.prefix(maxDisplayLength))
        }
        return clip.content
    }

    private var isTruncated: Bool {
        clip.content.count > maxDisplayLength
    }

    private var lineCount: Int {
        clip.content.components(separatedBy: .newlines).count
    }

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            header
            Divider()
            content
        }
        .glassEffect(.regular, in: .rect(cornerRadius: 12))
    }

    private var header: some View {
        VStack(alignment: .leading, spacing: 6) {
            HStack(spacing: 6) {
                Text(contentType)
                    .font(.caption2)
                    .fontWeight(.medium)
                    .padding(.horizontal, 8)
                    .padding(.vertical, 3)
                    .background(.quaternary)
                    .clipShape(Capsule())

                if clip.pinned {
                    Image(systemName: "pin.fill")
                        .font(.caption2)
                        .foregroundStyle(.orange)
                }

                Spacer()

                Button(action: onCopy) {
                    Label("Copy", systemImage: "doc.on.doc")
                        .font(.caption)
                }
                .buttonStyle(.plain)
                .foregroundStyle(.secondary)
            }

            HStack(spacing: 8) {
                Text(clip.date, style: .relative)
                    .font(.caption2)
                    .foregroundStyle(.secondary)

                Text("\(clip.content.count) chars")
                    .font(.caption2)
                    .foregroundStyle(.tertiary)

                Text("\(lineCount) lines")
                    .font(.caption2)
                    .foregroundStyle(.tertiary)
            }
        }
        .padding(10)
    }

    private var content: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 0) {
                if contentType == "Code" {
                    codeBlock
                } else {
                    plainTextBlock
                }

                if isTruncated {
                    Text("Content truncated...")
                        .font(.caption2)
                        .foregroundStyle(.tertiary)
                        .padding(.horizontal, 12)
                        .padding(.top, 4)
                        .padding(.bottom, 8)
                }
            }
        }
    }

    private var codeBlock: some View {
        Text(displayContent)
            .font(.system(size: 11, design: .monospaced))
            .foregroundStyle(.white)
            .textSelection(.enabled)
            .padding(12)
            .frame(maxWidth: .infinity, alignment: .leading)
            .background(Color(white: 0.12))
            .clipShape(RoundedRectangle(cornerRadius: 8))
            .padding(8)
    }

    private var plainTextBlock: some View {
        Text(displayContent)
            .font(.body)
            .textSelection(.enabled)
            .padding(12)
            .frame(maxWidth: .infinity, alignment: .leading)
    }
}
