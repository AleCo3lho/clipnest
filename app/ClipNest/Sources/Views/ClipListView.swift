import SwiftUI

struct ClipListView: View {
    @ObservedObject var viewModel: ClipboardViewModel

    private var pinnedClips: [Clip] {
        viewModel.clips.filter(\.pinned)
    }

    private var recentClips: [Clip] {
        viewModel.clips.filter { !$0.pinned }
    }

    var body: some View {
        HStack(spacing: 0) {
            if let hoveredClip = viewModel.hoveredClip {
                PreviewPanelView(clip: hoveredClip, onCopy: { viewModel.copyClip(hoveredClip) })
                    .frame(width: 308)
                    .padding(.leading, 6)
                    .padding(.vertical, 6)
                    .transition(.move(edge: .leading).combined(with: .opacity))
            } else {
                Spacer()
                    .frame(width: 314)
            }

            clipListContent
                .frame(width: 320)
        }
        .animation(.easeInOut(duration: 0.2), value: viewModel.hoveredClip?.id)
    }

    private var clipListContent: some View {
        VStack(spacing: 0) {
            // Search bar
            HStack(spacing: 8) {
                Image(systemName: "magnifyingglass")
                    .foregroundStyle(.secondary)
                TextField("Search clips...", text: $viewModel.searchQuery)
                    .textFieldStyle(.plain)
                    .onChange(of: viewModel.searchQuery) {
                        viewModel.search()
                    }
            }
            .padding(10)
            .glassEffect(.regular, in: .rect(cornerRadius: 10))
            .padding(.horizontal, 12)
            .padding(.top, 12)
            .padding(.bottom, 8)

            // Content
            if !viewModel.isConnected {
                StatusOverlayView(icon: "bolt.slash", message: "Daemon not running")
            } else if viewModel.clips.isEmpty {
                StatusOverlayView(icon: "clipboard", message: "No clips yet")
            } else {
                ScrollView {
                    LazyVStack(spacing: 6) {
                        if !pinnedClips.isEmpty {
                            sectionHeader("Pinned", icon: "pin.fill")
                            ForEach(pinnedClips) { clip in
                                ClipRowView(
                                    clip: clip,
                                    isHovered: viewModel.hoveredClip?.id == clip.id,
                                    onCopy: { viewModel.copyClip(clip) },
                                    onTogglePin: { viewModel.togglePin(clip) },
                                    onHover: { hovering in
                                        viewModel.setHoveredClip(hovering ? clip : nil)
                                    }
                                )
                            }
                        }

                        if !recentClips.isEmpty {
                            if !pinnedClips.isEmpty {
                                sectionHeader("Recent", icon: "clock")
                            }
                            ForEach(recentClips) { clip in
                                ClipRowView(
                                    clip: clip,
                                    isHovered: viewModel.hoveredClip?.id == clip.id,
                                    onCopy: { viewModel.copyClip(clip) },
                                    onTogglePin: { viewModel.togglePin(clip) },
                                    onHover: { hovering in
                                        viewModel.setHoveredClip(hovering ? clip : nil)
                                    }
                                )
                            }
                        }
                    }
                    .padding(.horizontal, 12)
                    .padding(.bottom, 8)
                }
            }

            Divider()

            // Footer
            HStack {
                Button("Clear All") {
                    viewModel.clearAll()
                }
                .buttonStyle(.plain)
                .foregroundStyle(.secondary)
                .font(.caption)

                Spacer()

                Button("Quit") {
                    NSApplication.shared.terminate(nil)
                }
                .buttonStyle(.plain)
                .foregroundStyle(.secondary)
                .font(.caption)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 10)
        }
    }

    private func sectionHeader(_ title: String, icon: String) -> some View {
        HStack(spacing: 4) {
            Image(systemName: icon)
                .font(.caption2)
            Text(title)
                .font(.caption)
                .fontWeight(.medium)
            Spacer()
        }
        .foregroundStyle(.secondary)
        .padding(.horizontal, 4)
        .padding(.top, 4)
    }
}
