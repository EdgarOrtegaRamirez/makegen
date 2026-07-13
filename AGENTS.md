# makegen — AGENTS.md

## Project Overview

makegen is a CLI tool that auto-generates Makefiles for Go, Rust, Python, and Node.js projects by detecting the project type and language.

## Architecture

```
makegen/
├── main.go                        # Entry point
├── cmd/
│   └── root.go                    # Cobra CLI commands and flags
├── internal/
│   ├── detect/
│   │   ├── detect.go              # Project type detection logic
│   │   └── detect_test.go         # 14 tests
│   └── generate/
│       ├── generate.go            # Makefile content generation
│       └── generate_test.go       # 7 tests
├── .github/workflows/ci.yml       # Go CI matrix (1.23, 1.24)
├── go.mod / go.sum
├── README.md
├── LICENSE
└── AGENTS.md
```

## Key Design Decisions

1. **Interface-free generation** — Each language has its own generation function (no interface overhead)
2. **String builders** — Uses `strings.Builder` for efficient Makefile construction
3. **Name extraction** — Custom parsers for go.mod, Cargo.toml, pyproject.toml, package.json (no heavy JSON/YAML parsing deps)
4. **Multi-type support** — Projects can have multiple languages (monorepos)
5. **No runtime dependencies beyond cobra** — Minimal dependency footprint
6. **Overwrite protection** — Requires `--force` flag to overwrite existing Makefiles

## Dependencies

- `github.com/spf13/cobra` — CLI framework

## Build & Test

```bash
go build -o makegen .
go test ./... -v
go vet ./...
```

## Adding a New Language

1. Add a `Type` constant in `internal/detect/detect.go`
2. Add detection logic in `DetectProject()` for the language's marker file
3. Add a generation function in `internal/generate/generate.go` (e.g., `pythonMakefile()`)
4. Add the function to the `generateForType()` switch
5. Add tests in both packages
6. Update README.md with the new language

## Common Tasks

### Fix name parsing for a project file
Edit the corresponding `parse*` function in `internal/detect/detect.go`:
- `parseGoModuleName()` — for go.mod
- `parseCargoName()` — for Cargo.toml
- `parsePyprojectName()` — for pyproject.toml
- `parsePackageName()` — for package.json

### Add a new Makefile target
1. Add the target definition in the language-specific generation function
2. Add it to the `.PHONY` line and help section
3. Update tests

## Troubleshooting

- **Tests fail with "redeclared"**: Check for naming conflicts between type names and function names
- **Name detection returning "."**: The parser for that language isn't matching — check the `parse*` function