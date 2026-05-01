# Diary

Diary is a local-first CLI for AI harness implementation memory. It records compact handoffs into `.diary/`, retrieves prompt-ready context for the next run, and installs optional harness skills.

## Install

Install the latest release:

```bash
curl -fsSL https://raw.githubusercontent.com/webteractive/diary/main/scripts/install.sh | sh
```

Install a specific version:

```bash
curl -fsSL https://raw.githubusercontent.com/webteractive/diary/main/scripts/install.sh | env DIARY_VERSION=v0.0.1 sh
```

Install somewhere other than `/usr/local/bin`:

```bash
curl -fsSL https://raw.githubusercontent.com/webteractive/diary/main/scripts/install.sh | env DIARY_INSTALL_DIR="$HOME/.local/bin" sh
```

## Build From Source

```bash
go build -o bin/diary ./cmd/diary
```

## Usage

Record compact implementation context:

```bash
diary record "Implemented record storage and added tests."
```

Retrieve the latest context:

```bash
diary get
```

List stored records or projects:

```bash
diary list
diary list --projects
```

Install skills for supported harness targets:

```bash
diary install-skills --target codex
diary install-skills --target claude
```

This installs `diary-get`, `diary-record`, and `diary-list`.

Update Diary from GitHub releases:

```bash
diary self-update
diary self-update --version v0.0.1
diary self-update --dry-run
```

## Development

```bash
go test ./...
go run ./cmd/diary --help
```
