# Homebrew formula for ClipNest
class Clipnest < Formula
  desc "Your cozy clipboard manager for macOS"
  homepage "https://github.com/AleCo3lho/clipnest"
  url "https://github.com/AleCo3lho/clipnest/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  license "MIT"

  depends_on "go" => :build
  depends_on "sqlite"

  def install
    system "go", "build", "-o", "clipnest", "./cmd/clipnest"
    system "go", "build", "-o", "clipnestd", "./cmd/clipnestd"
    bin.install "clipnest"
    bin.install "clipnestd"

    # Create support directory
    (var/"clipnest").mkpath
  end

  test do
    system "clipnest", "version"
  end
end
