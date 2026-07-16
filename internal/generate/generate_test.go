package generate

import (
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/makegen/internal/detect"
)

func checkContains(t *testing.T, content, needle string) {
	t.Helper()
	if !strings.Contains(content, needle) {
		t.Errorf("expected content to contain %q", needle)
	}
}

func TestGoMakefile(t *testing.T) {
	project := &detect.Project{
		Type:     detect.TypeGo,
		Name:     "github.com/test/myapp",
		Dir:      ".",
		SubTypes: []detect.Type{detect.TypeGo},
		Metadata: map[string]string{"go_version": "1.24.0"},
	}

	content, err := Makefile(project, "myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checkContains(t, content, "PROJECT_NAME ?= myapp")
	checkContains(t, content, "build-go")
	checkContains(t, content, "test-go")
	checkContains(t, content, "clean-go")
	checkContains(t, content, "lint-go")
	checkContains(t, content, "format-go")
	checkContains(t, content, "go build")
	checkContains(t, content, "go test")
	checkContains(t, content, "go vet")
	checkContains(t, content, "go fmt")
	checkContains(t, content, "build: build-go")
	checkContains(t, content, "GO_VERSION ?= 1.24.0")
}

func TestRustMakefile(t *testing.T) {
	project := &detect.Project{
		Type:     detect.TypeRust,
		Name:     "myapp",
		Dir:      ".",
		SubTypes: []detect.Type{detect.TypeRust},
	}

	content, err := Makefile(project, "myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checkContains(t, content, "PROJECT_NAME ?= myapp")
	checkContains(t, content, "build-rust")
	checkContains(t, content, "test-rust")
	checkContains(t, content, "clean-rust")
	checkContains(t, content, "cargo build")
	checkContains(t, content, "cargo test")
	checkContains(t, content, "cargo clean")
}

func TestPythonMakefile(t *testing.T) {
	project := &detect.Project{
		Type:     detect.TypePython,
		Name:     "mypkg",
		Dir:      ".",
		SubTypes: []detect.Type{detect.TypePython},
	}

	content, err := Makefile(project, "mypkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checkContains(t, content, "PROJECT_NAME ?= mypkg")
	checkContains(t, content, "build-python")
	checkContains(t, content, "test-python")
	checkContains(t, content, "clean-python")
	checkContains(t, content, "$(PYTHON) -m build")
	checkContains(t, content, "$(PYTHON) -m pytest")
	checkContains(t, content, "ruff")
}

func TestNodeMakefile(t *testing.T) {
	project := &detect.Project{
		Type:     detect.TypeNode,
		Name:     "my-app",
		Dir:      ".",
		SubTypes: []detect.Type{detect.TypeNode},
	}

	content, err := Makefile(project, "my-app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checkContains(t, content, "PROJECT_NAME ?= my-app")
	checkContains(t, content, "build-node")
	checkContains(t, content, "test-node")
	checkContains(t, content, "clean-node")
	checkContains(t, content, "$(NPM) run build")
	checkContains(t, content, "$(NPM) test")
	checkContains(t, content, "node_modules")
}

func TestMultiMakefile(t *testing.T) {
	project := &detect.Project{
		Type:     detect.TypeMulti,
		Name:     "monorepo",
		Dir:      ".",
		SubTypes: []detect.Type{detect.TypeGo, detect.TypeNode},
	}

	content, err := Makefile(project, "monorepo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checkContains(t, content, "PROJECT_NAME ?= monorepo")
	checkContains(t, content, "build-go")
	checkContains(t, content, "build-node")
	checkContains(t, content, "build: build-go build-node")
	checkContains(t, content, "test: test-go test-node")
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"simple", "simple"},
		{"with-dashes", "with-dashes"},
		{"github.com/user/repo", "repo"},
		{"special!@#chars", "special___chars"},
		{"UPPERCASE", "UPPERCASE"},
	}
	for _, tt := range tests {
		got := sanitizeName(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestHelpTarget(t *testing.T) {
	project := &detect.Project{
		Type:     detect.TypeGo,
		Name:     "test",
		Dir:      ".",
		SubTypes: []detect.Type{detect.TypeGo},
	}

	content, err := Makefile(project, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checkContains(t, content, ".PHONY:")
	checkContains(t, content, "help:")
	checkContains(t, content, "release:")
	checkContains(t, content, "install:")
	checkContains(t, content, "audit:")
}
