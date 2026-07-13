# makegen 🛠️

[![CI](https://github.com/EdgarOrtegaRamirez/makegen/actions/workflows/ci.yml/badge.svg)](https://github.com/EdgarOrtegaRamirez/makegen/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/EdgarOrtegaRamirez/makegen)](https://goreportcard.com/report/github.com/EdgarOrtegaRamirez/makegen)

**Auto-generate Makefiles for Go, Rust, Python, and Node.js projects.**

Makegen detects the project type in your directory and generates a comprehensive Makefile with common targets — build, test, clean, lint, format, release, install, and audit.

```bash
# Generate a Makefile for the current directory
makegen

# Generate for a specific directory
makegen /path/to/project

# Overwrite existing Makefile
makegen --force

# List supported project types
makegen --list-types
```

## Features

- **🔍 Auto-detection** — Detects project type by scanning for `go.mod`, `Cargo.toml`, `pyproject.toml`, `package.json`, etc.
- **🦀 4 languages** — Go, Rust, Python, Node.js
- **📦 Multi-type support** — Projects with multiple languages (e.g., Go + Node.js monorepos)
- **🎯 Per-language targets** — Language-specific build, test, clean, lint, and format
- **🔗 Unified targets** — `build`, `test`, `clean`, `lint`, `format` that delegate to the right language
- **🧹 Smart defaults** — Sensible build flags, test flags, and cleanup patterns
- **📝 Project name extraction** — Reads name from `go.mod`, `Cargo.toml`, `pyproject.toml`, or `package.json`
- **🏷️ Custom naming** — Override project name with `--name`

## Installation

```bash
# From source
go install github.com/EdgarOrtegaRamirez/makegen@latest

# Or build from repo
git clone https://github.com/EdgarOrtegaRamirez/makegen.git
cd makegen
go build -o makegen .
sudo mv makegen /usr/local/bin/
```

## Usage

### Basic

```bash
# Generate Makefile in current directory
$ makegen
Generated Makefile for go project (myapp)
Output: ./Makefile

# Generate for a specific directory
$ makegen /path/to/rust/project
Generated Makefile for rust project (myapp)
Output: /path/to/rust/project/Makefile

# Detect project name automatically
$ makegen /path/to/project
Generated Makefile for go project (github.com/user/project)
Output: /path/to/project/Makefile
```

### Options

```bash
# Override output path
makegen --output /tmp/Makefile

# Override project name
makegen --name my-cool-app

# Force overwrite existing Makefile
makegen --force

# List supported project types
makegen --list-types
```

### Generated Makefile targets

| Target    | Description                            |
|-----------|----------------------------------------|
| `build`   | Build the project                      |
| `test`    | Run tests                              |
| `clean`   | Clean build artifacts                  |
| `lint`    | Run linters (go vet, cargo check, etc) |
| `format`  | Format code                            |
| `release` | Build release artifacts                |
| `install` | Install the project                    |
| `audit`   | Run security audit                     |
| `help`    | Show all targets                       |

Each language also has language-specific targets:
- `build-go`, `test-go`, `clean-go`, `lint-go`, `format-go`
- `build-rust`, `test-rust`, `clean-rust`, `lint-rust`, `format-rust`
- `build-python`, `test-python`, `clean-python`, `lint-python`, `format-python`
- `build-node`, `test-node`, `clean-node`, `lint-node`, `format-node`

## Supported Project Types

| Type     | Detection Files                    | Languages/Scripts          |
|----------|------------------------------------|----------------------------|
| Go       | `go.mod`                           | `go`                       |
| Rust     | `Cargo.toml`                       | `cargo`, `rustfmt`         |
| Python   | `pyproject.toml`, `setup.py`, etc  | `python3`, `pip`, `pytest`, `ruff`, `mypy` |
| Node.js  | `package.json`                     | `npm`, `npx`, `prettier`   |
| Multi    | Multiple detection files present   | All detected languages     |

## Development

```bash
# Build
go build -o makegen .

# Run all tests
go test ./... -v

# Lint
go vet ./...

# Format
go fmt ./...
```

## How It Works

1. **Detect** — Scans the target directory for well-known project files
2. **Analyze** — Extracts project metadata (name, Go version, etc.)
3. **Generate** — Creates a Makefile with appropriate targets for each detected language
4. **Write** — Saves the Makefile (checks for existing file, requires `--force` to overwrite)

## License

MIT