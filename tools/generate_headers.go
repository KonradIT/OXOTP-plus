// Generate gzipped C header files from web assets.
//
// Usage:
//
//	go run tools/generate_headers.go
//
// This script converts the HTML/CSS/favicon files to gzipped C arrays
// that can be served by the ESP32 web server.
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const bytesPerLine = 22

func main() {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("OXOTP+ Web Asset Header Generator (Go)")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	// Get project root (parent of tools directory)
	exePath, err := os.Executable()
	if err != nil {
		// Fallback: use current working directory
		exePath, _ = os.Getwd()
	}

	// Try to find project root by looking for platformio.ini
	projectRoot := findProjectRoot(exePath)
	if projectRoot == "" {
		// Fallback to current directory
		projectRoot, _ = os.Getwd()
	}

	webDir := filepath.Join(projectRoot, "web")
	includeDir := filepath.Join(projectRoot, "include")

	// Ensure include directory exists
	if err := os.MkdirAll(includeDir, 0755); err != nil {
		fmt.Printf("Error creating include directory: %v\n", err)
		os.Exit(1)
	}

	// Generate index.h
	if err := generateHeader(
		filepath.Join(webDir, "index.html"),
		filepath.Join(includeDir, "index.h"),
		"index_html_gz",
	); err != nil {
		fmt.Printf("Error generating index.h: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()

	// Generate favico.h
	if err := generateHeader(
		filepath.Join(webDir, "favicon.png"),
		filepath.Join(includeDir, "favico.h"),
		"favico_gz",
	); err != nil {
		fmt.Printf("Error generating favico.h: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Done! Headers have been regenerated.")
	fmt.Println("Rebuild the firmware to apply changes.")
}

// findProjectRoot walks up directories looking for platformio.ini
func findProjectRoot(startPath string) string {
	dir := startPath
	for i := 0; i < 10; i++ { // Max 10 levels up
		if _, err := os.Stat(filepath.Join(dir, "platformio.ini")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// generateHeader reads a file, gzips it, and writes a C header
func generateHeader(inputPath, outputPath, varName string) error {
	fmt.Printf("Processing: %s\n", inputPath)

	// Read input file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Gzip compress
	var buf bytes.Buffer
	gzWriter, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return fmt.Errorf("failed to create gzip writer: %w", err)
	}

	if _, err := gzWriter.Write(data); err != nil {
		return fmt.Errorf("failed to write gzip data: %w", err)
	}

	if err := gzWriter.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}

	compressed := buf.Bytes()

	// Generate C array
	cCode := formatCArray(compressed, varName)

	// Write output file
	if err := os.WriteFile(outputPath, []byte(cCode), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("  Generated: %s (%d bytes -> %d bytes gzipped)\n",
		outputPath, len(data), len(compressed))

	return nil
}

// formatCArray converts bytes to a C array string
func formatCArray(data []byte, varName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("const char %s[] PROGMEM = {\n", varName))

	for i := 0; i < len(data); i += bytesPerLine {
		end := i + bytesPerLine
		if end > len(data) {
			end = len(data)
		}

		chunk := data[i:end]
		hexValues := make([]string, len(chunk))
		for j, b := range chunk {
			hexValues[j] = fmt.Sprintf("0x%02x", b)
		}

		line := strings.Join(hexValues, ",")
		if end < len(data) {
			line += ","
		}
		sb.WriteString(line + "\n")
	}

	sb.WriteString("};")

	return sb.String()
}
