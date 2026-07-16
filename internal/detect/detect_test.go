package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectGo(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module github.com/test/project\n\ngo 1.24.0\n"), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := DetectProject(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Type != TypeGo {
		t.Errorf("expected Go, got %s", p.Type)
	}
	if p.Name != "github.com/test/project" {
		t.Errorf("expected github.com/test/project, got %s", p.Name)
	}
}

func TestDetectRust(t *testing.T) {
	dir := t.TempDir()
	content := `[package]
name = "myapp"
version = "0.1.0"
edition = "2024"
`
	if err := os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := DetectProject(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Type != TypeRust {
		t.Errorf("expected Rust, got %s", p.Type)
	}
	if p.Name != "myapp" {
		t.Errorf("expected myapp, got %s", p.Name)
	}
}

func TestDetectPython(t *testing.T) {
	dir := t.TempDir()
	content := `[project]
name = "mypackage"
version = "0.1.0"
`
	if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := DetectProject(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Type != TypePython {
		t.Errorf("expected Python, got %s", p.Type)
	}
	if p.Name != "mypackage" {
		t.Errorf("expected mypackage, got %s", p.Name)
	}
}

func TestDetectNode(t *testing.T) {
	dir := t.TempDir()
	content := `{
  "name": "my-node-app",
  "version": "1.0.0"
}
`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := DetectProject(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Type != TypeNode {
		t.Errorf("expected Node, got %s", p.Type)
	}
	if p.Name != "my-node-app" {
		t.Errorf("expected my-node-app, got %s", p.Name)
	}
}

func TestDetectRequirements(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "requirements.txt"), []byte("requests\nflask\n"), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := DetectProject(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Type != TypePython {
		t.Errorf("expected Python, got %s", p.Type)
	}
	// Name should fall back to directory name
	if p.Name != filepath.Base(dir) {
		t.Errorf("expected %s, got %s", filepath.Base(dir), p.Name)
	}
}

func TestDetectUnknown(t *testing.T) {
	dir := t.TempDir()
	// Empty directory
	if err := os.WriteFile(filepath.Join(dir, "random.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := DetectProject(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Type != TypeUnknown {
		t.Errorf("expected Unknown, got %s", p.Type)
	}
}

func TestDetectMulti(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module multi\n\ngo 1.24.0\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name":"multi-node"}`), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := DetectProject(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Type != TypeMulti {
		t.Errorf("expected Multi, got %s", p.Type)
	}
	if len(p.SubTypes) != 2 {
		t.Errorf("expected 2 subtypes, got %d", len(p.SubTypes))
	}
}

func TestDetectNotADir(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file")
	if err := os.WriteFile(f, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := DetectProject(f)
	if err == nil {
		t.Fatal("expected error for non-directory")
	}
}

func TestDetectNoDir(t *testing.T) {
	_, err := DetectProject("/nonexistent/path/xyz123")
	if err == nil {
		t.Fatal("expected error for nonexistent path")
	}
}

func TestParseGoVersion(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"module foo\n\ngo 1.24.0\n", "1.24.0"},
		{"go 1.21", "1.21"},
		{"module foo\n", ""},
	}
	for _, tt := range tests {
		got := parseGoVersion(tt.input)
		if got != tt.want {
			t.Errorf("parseGoVersion(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseCargoName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`[package]
name = "myapp"`, "myapp"},
		{`[package]
name = "cool-tool"`, "cool-tool"},
		{``, ""},
	}
	for _, tt := range tests {
		got := parseCargoName(tt.input)
		if got != tt.want {
			t.Errorf("parseCargoName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParsePyprojectName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`[project]
name = "mypkg"
version = "0.1.0"`, "mypkg"},
		{`[build-system]
requires = ["setuptools"]`, ""},
		{``, ""},
	}
	for _, tt := range tests {
		got := parsePyprojectName(tt.input)
		if got != tt.want {
			t.Errorf("parsePyprojectName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParsePackageName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`{"name": "my-app"}`, "my-app"},
		{`{"name": "test-pkg","version":"1.0"}`, "test-pkg"},
		{`{}`, ""},
	}
	for _, tt := range tests {
		got := parsePackageName(tt.input)
		if got != tt.want {
			t.Errorf("parsePackageName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGoVersion(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example\n\ngo 1.23.5\n"), 0644); err != nil {
		t.Fatal(err)
	}
	p, err := DetectProject(dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.Metadata["go_version"] != "1.23.5" {
		t.Errorf("expected go_version=1.23.5, got %s", p.Metadata["go_version"])
	}
}
