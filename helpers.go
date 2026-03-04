package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func directoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		// Path exists, let's verify it's actually a directory
		return info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		// Path explicitly does not exist
		return false, nil
	}
	// An error occurred (e.g., permission issues)
	return false, err
}

func setupTrashDir(path string) (string, error) {
	// Ensure the root trash path exists and has strict permissions
	if err := os.MkdirAll(path, 0o700); err != nil {
		return "", fmt.Errorf("failed to create trash root %s: %w", path, err)
	}
	// Explicitly enforce permissions in case the dir already exixted with
	// loose permissions (eg: 0o755)
	if err := os.Chmod(path, 0o700); err != nil {
		return "", fmt.Errorf("failed to secure trash directory: %w", err)
	}

	subdirs := []string{
		filepath.Join(path, "files"), // actual files
		filepath.Join(path, "info"),  // metadata files
	}
	for _, dir := range subdirs {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return "", fmt.Errorf(
				"failed to create trash subdir %s: %w", dir, err,
			)
		}
	}
	dirCacheFile := filepath.Join(path, "directorysizes")
	if err := makeDirCacheFile(dirCacheFile); err != nil {
		return "", err
	}
	return path, nil
}

func makeDirCacheFile(cachepath string) error {
	// O_CREATE | O_EXCL ensures we don't truncate or touch it if it exists
	f, err := os.OpenFile(cachepath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil
		}
		return fmt.Errorf("could not create cache file %s: %w", cachepath, err)
	}
	return f.Close()
}

// fileCheck checks the file existence, and return its absolute path
// if exist otherwise an error.
func fileCheck(file string) (string, error) {
	if _, err := os.Stat(file); err != nil {
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return "", fmt.Errorf("File is missing: %w", err)
		case errors.Is(err, fs.ErrPermission):
			return "", fmt.Errorf("Permission Denied for %s: %w", file, err)
		default:
			return "", fmt.Errorf("System Error during Stat: %w", err)
		}
	}
	absPath, err := filepath.Abs(file)
	if err != nil {
		return "", fmt.Errorf("Could Not determine absolute path: %w", err)
	}
	return absPath, nil
}
