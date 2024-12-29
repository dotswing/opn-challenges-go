package fileutils

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadFileToBytes(filePath string) ([]byte, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Printf("Error getting absolute path: %v\n", err)
		return nil, err
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return nil, err
	}
	return data, nil
}
