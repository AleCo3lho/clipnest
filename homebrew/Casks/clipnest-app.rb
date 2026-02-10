cask "clipnest-app" do
  version "PLACEHOLDER"
  sha256 "PLACEHOLDER"

  url "https://github.com/AleCo3lho/clipnest/releases/download/v#{version}/ClipNest.app.zip"
  name "ClipNest"
  desc "Privacy-first clipboard manager for macOS"
  homepage "https://github.com/AleCo3lho/clipnest"

  depends_on formula: "AleCo3lho/clipnest/clipnest"
  depends_on macos: ">= :tahoe"

  app "ClipNest.app"

  postflight do
    system_command "/opt/homebrew/bin/brew",
                   args: ["services", "start", "AleCo3lho/clipnest/clipnest"],
                   sudo: false
  end

  uninstall quit: "com.clipnest.app"

  zap trash: [
    "/tmp/clipnest.sock",
  ]
end
