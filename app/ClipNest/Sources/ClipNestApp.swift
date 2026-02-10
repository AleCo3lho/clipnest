import SwiftUI

@main
struct ClipNestApp: App {
    @StateObject private var viewModel = ClipboardViewModel()

    var body: some Scene {
        MenuBarExtra {
            ClipListView(viewModel: viewModel)
                .frame(width: 640, height: 440)
        } label: {
            Image(systemName: "clipboard")
        }
        .menuBarExtraStyle(.window)
    }
}
