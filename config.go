package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// loadRiddle reads a file named "riddle.txt" in the provided config directory.
// This file should have a bcrypt 12 cost hash of the riddle's answer on the
// first line, followed by the riddle's body in plain text for the remainder of
// the file.
func loadRiddle(config string) (string, string, error) {
	data, err := os.ReadFile(filepath.Join(config, "riddle.txt"))
	if err != nil {
		return "", "", err
	}

	before, after, ok := strings.Cut(string(data), "\n")
	if !ok {
		return "", "", fmt.Errorf("riddle.txt does not contain a newline\n")
	}
	return after, before, nil
}
