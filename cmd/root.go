package cmd

import (
	"fmt"
	"os"

	"github.com/EdgarOrtegaRamirez/makegen/internal/detect"
	"github.com/EdgarOrtegaRamirez/makegen/internal/generate"
	"github.com/spf13/cobra"
)

var (
	outputFile     string
	projectName    string
	forceOverwrite bool
	listTypes      bool
)

var rootCmd = &cobra.Command{
	Use:   "makegen [directory]",
	Short: "Auto-generate Makefiles for Go, Rust, Python, and Node.js projects",
	Long: `makegen detects the project type in the given directory (or current directory)
and generates a comprehensive Makefile with common targets (build, test, clean, 
lint, format, release, install, and more).

Supported project types:
  - Go       (go.mod)
  - Rust     (Cargo.toml)
  - Python   (pyproject.toml, setup.py, setup.cfg, requirements.txt)
  - Node.js  (package.json)
  - Multi    (detects multiple project types in one directory)`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if listTypes {
			printSupportedTypes()
			return nil
		}

		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		// Detect project type
		project, err := detect.DetectProject(dir)
		if err != nil {
			return fmt.Errorf("detect: %w", err)
		}

		if project.Type == detect.TypeUnknown {
			return fmt.Errorf("no supported project type detected in %s", dir)
		}

		// Determine project name
		name := projectName
		if name == "" {
			name = project.Name
		}
		if name == "" {
			name = "project"
		}

		// Generate Makefile content
		content, err := generate.Makefile(project, name)
		if err != nil {
			return fmt.Errorf("generate: %w", err)
		}

		// Write output
		outPath := outputFile
		if outPath == "" {
			outPath = dir + "/Makefile"
		}

		// Check existing file
		if _, err := os.Stat(outPath); err == nil && !forceOverwrite {
			return fmt.Errorf("%s already exists (use --force to overwrite)", outPath)
		}

		if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("write: %w", err)
		}

		fmt.Printf("Generated Makefile for %s project (%s)\n", project.Type, name)
		fmt.Printf("Output: %s\n", outPath)
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output path (default: <dir>/Makefile)")
	rootCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name (auto-detected if not set)")
	rootCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "Overwrite existing Makefile")
	rootCmd.Flags().BoolVarP(&listTypes, "list-types", "l", false, "List supported project types")
}

func printSupportedTypes() {
	fmt.Println("Supported project types:")
	fmt.Println("  go       - Go projects (go.mod)")
	fmt.Println("  rust     - Rust projects (Cargo.toml)")
	fmt.Println("  python   - Python projects (pyproject.toml, setup.py, etc.)")
	fmt.Println("  node     - Node.js projects (package.json)")
	fmt.Println("  multi    - Multiple project types in one directory")
}
