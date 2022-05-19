package fileio

import (
	"fmt"
	"os"
	"path/filepath"
)

func PrepareWDir() (string, error) {
	dir := workingDir()
	if err := prepareDir(dir); err != nil {
		return "", err
	}
	return dir, nil
}

func workingDir() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".notes")
}

func prepareDir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat: %s: %w", dir, err)
		}
		// Dir Not Exsist
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("mkdir_all: %s: %w", dir, err)
		}
	}
	return nil // Dir Exists
}
