# ClipNest Homebrew Tap

This directory contains the source-of-truth templates for the [homebrew-clipnest](https://github.com/AleCo3lho/homebrew-clipnest) tap. The release pipeline copies these files to the tap repo and updates version/SHA256 values automatically.

## Structure

```
homebrew/
  Formula/clipnest.rb     # Builds clipnest + clipnestd from Go source
  Casks/clipnest-app.rb    # Downloads pre-built ClipNest.app
```

## Install

Install the GUI app (includes CLI + daemon):

```
brew install --cask AleCo3lho/clipnest/clipnest-app
```

Install just the CLI and daemon (no GUI):

```
brew install AleCo3lho/clipnest/clipnest
brew services start clipnest
```

## How releases work

1. Push a tag `v*` to trigger the release workflow
2. The `build` job (macOS) builds Go binaries + Swift app, creates a GitHub release
3. The `update-tap` job downloads assets, computes SHA256 hashes, and pushes updates to the tap repo

## Setup (one-time)

The release pipeline requires a `TAP_GITHUB_TOKEN` secret in this repo with write access to `AleCo3lho/homebrew-clipnest`.
