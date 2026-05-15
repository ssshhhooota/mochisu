# ghissues

A terminal UI for browsing GitHub issues.

## Features

- List repositories for a user or organization
- Browse issues in a selected repository
- Open issues in your web browser

## Installation

### Homebrew

```bash
brew install ssshhhooota/tap/ghissues
```

### Build from source

```bash
git clone https://github.com/ssshhhooota/ghissues
cd ghissues
make build
```

## Requirements

- [GitHub CLI (`gh`)](https://cli.github.com/) authenticated via `gh auth login`

## Usage

```bash
# List your own repositories
ghissues

# List repositories for a specific user or organization
ghissues <owner>
```

In the TUI, select a repository to view its issues, then select an issue to open it in your browser.

## Development

```bash
make setup     # Install dependencies and tools
make dev       # Run with hot reload (requires gow)
make run       # Run once
make lint      # Run golangci-lint
make build     # Build binary into ./bin/
```

## Release

See [RELEASE.md](./RELEASE.md).

## License

[MIT](./LICENSE)
