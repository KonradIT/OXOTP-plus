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
	fmt.Println("OXOTP+ Web Asset Header Generator")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	webDir := filepath.Join(".", "web")
	includeDir := filepath.Join(".", "include")

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

	_, err = gzWriter.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write gzip data: %w", err)
	}

	err = gzWriter.Close()
	if err != nil {
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
		end := min(i+bytesPerLine, len(data))

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
