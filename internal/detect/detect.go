package detect

import (
	"fmt"
	"os"
	"path/filepath"
)

// Type represents the project type.
type Type string

const (
	TypeGo     Type = "go"
	TypeRust   Type = "rust"
	TypePython Type = "python"
	TypeNode   Type = "node"
	TypeMulti  Type = "multi"
	TypeUnknown Type = "unknown"
)

// Project describes a detected project.
type Project struct {
	Type Type
	Name string
	Dir  string
	Files []string          // Files that triggered detection
	SubTypes []Type         // For multi-type projects
	Metadata map[string]string // Extra info (Go version, Rust edition, etc.)
}

// DetectionFiles returns common detection files for a type.
func (t Type) DetectionFiles() []string {
	switch t {
	case TypeGo:
		return []string{"go.mod"}
	case TypeRust:
		return []string{"Cargo.toml"}
	case TypePython:
		return []string{"pyproject.toml", "setup.py", "setup.cfg", "requirements.txt"}
	case TypeNode:
		return []string{"package.json"}
	default:
		return nil
	}
}

// DetectProject detects the project type in the given directory.
func DetectProject(dir string) (*Project, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("accessing %s: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}

	var detected []Type
	fileMap := make(map[string]bool)
	for _, e := range entries {
		if !e.IsDir() {
			fileMap[e.Name()] = true
		}
	}

	var files []string
	var subTypes []Type
	metadata := make(map[string]string)

	if fileMap["go.mod"] {
		detected = append(detected, TypeGo)
		files = append(files, "go.mod")
		subTypes = append(subTypes, TypeGo)
		// Read Go version from go.mod
		if data, err := os.ReadFile(filepath.Join(dir, "go.mod")); err == nil {
			metadata["go_version"] = parseGoVersion(string(data))
		}
	}

	if fileMap["Cargo.toml"] {
		detected = append(detected, TypeRust)
		files = append(files, "Cargo.toml")
		subTypes = append(subTypes, TypeRust)
		// Read package name from Cargo.toml
		if data, err := os.ReadFile(filepath.Join(dir, "Cargo.toml")); err == nil {
			metadata["cargo_name"] = parseCargoName(string(data))
		}
	}

	if fileMap["pyproject.toml"] || fileMap["setup.py"] || fileMap["setup.cfg"] || fileMap["requirements.txt"] {
		detected = append(detected, TypePython)
		subTypes = append(subTypes, TypePython)
		if fileMap["pyproject.toml"] {
			files = append(files, "pyproject.toml")
		} else if fileMap["setup.py"] {
			files = append(files, "setup.py")
		} else if fileMap["setup.cfg"] {
			files = append(files, "setup.cfg")
		} else {
			files = append(files, "requirements.txt")
		}
	}

	if fileMap["package.json"] {
		detected = append(detected, TypeNode)
		files = append(files, "package.json")
		subTypes = append(subTypes, TypeNode)
	}

	if len(detected) == 0 {
		return &Project{Type: TypeUnknown, Dir: dir}, nil
	}

	projectType := detected[0]
	if len(detected) > 1 {
		projectType = TypeMulti
	}

	// Extract project name
	name := detectName(dir, detected, fileMap)

	return &Project{
		Type:     projectType,
		Name:     name,
		Dir:      dir,
		Files:    files,
		SubTypes: subTypes,
		Metadata: metadata,
	}, nil
}

func detectName(dir string, types []Type, files map[string]bool) string {
	// Try to get name from project files
	for _, t := range types {
		switch t {
		case TypeGo:
			if data, err := os.ReadFile(filepath.Join(dir, "go.mod")); err == nil {
				if n := parseGoModuleName(string(data)); n != "" {
					return n
				}
			}
		case TypeRust:
			if data, err := os.ReadFile(filepath.Join(dir, "Cargo.toml")); err == nil {
				if n := parseCargoName(string(data)); n != "" {
					return n
				}
			}
		case TypePython:
			if data, err := os.ReadFile(filepath.Join(dir, "pyproject.toml")); err == nil {
				if n := parsePyprojectName(string(data)); n != "" {
					return n
				}
			}
		case TypeNode:
			// Read from package.json
			if data, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
				if n := parsePackageName(string(data)); n != "" {
					return n
				}
			}
		}
	}

	// Fallback: use directory name
	return filepath.Base(dir)
}

func parseGoVersion(content string) string {
	// Simple parser for: go 1.xx
	for _, line := range splitLines(content) {
		line = trimSpace(line)
		if len(line) > 3 && line[:3] == "go " && line[3] >= '1' && line[3] <= '9' {
			return line[3:]
		}
	}
	return ""
}

func parseGoModuleName(content string) string {
	// First line is typically: module <name>
	for _, line := range splitLines(content) {
		line = trimSpace(line)
		if len(line) > 7 && line[:7] == "module " {
			return line[7:]
		}
	}
	return ""
}

func parseCargoName(content string) string {
	for _, line := range splitLines(content) {
		line = trimSpace(line)
		if len(line) > 7 && line[:7] == "name = " {
			// name = "foo" or name = 'foo'
			val := line[7:]
			if len(val) >= 2 {
				val = val[1 : len(val)-1]
			}
			return val
		}
	}
	return ""
}

func parsePyprojectName(content string) string {
	inProject := false
	for _, line := range splitLines(content) {
		trimmed := trimSpace(line)
		if trimmed == "[project]" {
			inProject = true
			continue
		}
		if inProject {
			if len(trimmed) > 0 && trimmed[0] == '[' {
				break // next section
			}
			if len(trimmed) > 7 && trimmed[:7] == "name = " {
				val := trimmed[7:]
				if len(val) >= 2 {
					val = val[1 : len(val)-1]
				}
				return val
			}
		}
	}
	return ""
}
func parsePackageName(content string) string {
	// Simple JSON parse for "name" field
	for _, line := range splitLines(content) {
		trimmed := trimSpace(line)
		// Look for "name": " anywhere in the line
		for i := 0; i+9 <= len(trimmed); i++ {
			if trimmed[i:i+9] == `"name": "` {
				start := i + 9
				end := start
				for end < len(trimmed) && trimmed[end] != '"' {
					end++
				}
				if end > start {
					return trimmed[start:end]
				}
			}
		}
	}
	return ""
}

// splitLines splits a string into lines.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// trimSpace trims spaces and tabs.
func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}