package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetCurrentWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("Error getting working directory: %v", err)
	}
	return dir
}

func GetAbsolutePath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Sprintf("Error getting absolute path: %v", err)
	}
	return abs
}
