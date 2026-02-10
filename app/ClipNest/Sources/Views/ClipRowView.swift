import SwiftUI

struct ClipRowView: View {
    let clip: Clip
    let isHovered: Bool
    let now: Date
    let onCopy: () -> Void
    let onTogglePin: () -> Void
    let onHover: (Bool) -> Void

    var body: some View {
        Button(action: onCopy) {
            VStack(alignment: .leading, spacing: 4) {
                Text(clip.preview)
                    .font(.system(.body, design: .monospaced))
                    .lineLimit(2)
                    .frame(maxWidth: .infinity, alignment: .leading)

                HStack {
                    Text(ContentClassifier.classify(clip.content))
                        .font(.caption2)
                        .fontWeight(.medium)
                        .padding(.horizontal, 6)
                        .padding(.vertical, 2)
                        .background(.quaternary)
                        .clipShape(Capsule())

                    Text(RelativeTimeFormatter.string(for: clip.date, relativeTo: now))
                        .font(.caption)
                        .foregroundStyle(.secondary)

                    Spacer()

                    if clip.pinned {
                        Image(systemName: "pin.fill")
                            .font(.caption)
                            .foregroundStyle(.orange)
                    }
                }
            }
            .padding(10)
            .contentShape(Rectangle())
        }
        .buttonStyle(.plain)
        .modifier(GlassBackgroundModifier(cornerRadius: 12))
        .overlay(
            RoundedRectangle(cornerRadius: 12)
                .strokeBorder(.primary.opacity(isHovered ? 0.2 : 0), lineWidth: 1)
        )
        .onHover { hovering in onHover(hovering) }
        .contextMenu {
            Button {
                onCopy()
            } label: {
                Label("Copy", systemImage: "doc.on.doc")
            }
            Button {
                onTogglePin()
            } label: {
                Label(clip.pinned ? "Unpin" : "Pin", systemImage: clip.pinned ? "pin.slash" : "pin")
            }
        }
    }
}
