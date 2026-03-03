package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type TrashImpl interface {
	TopDir(mountPoint string) (string, error)
	TrashDir() (string, error)
}

type HomeTrash struct {
	topDir string
}

func (ht *HomeTrash) TopDir(mountPoint string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not locate home directory: %w", err)
	}
	return home, nil
}

func (ht *HomeTrash) TrashDir() (string, error) {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := ht.TopDir()
		if err != nil {
			return "", err
		}
		dataHome = filepath.Join(home, ".local", "share")
	}
	trashPath := filepath.Join(dataHome, "Trash")
	return setupTrashDir(trashPath)
}

type OtherTrash struct {
	topDir string
}

func (ot *OtherTrash) TopDir(mountPoint string) (string, error) {
	// TODO:
	return mountPoint, nil
}

func (ot *OtherTrash) TrashDir() (string, error) {
	// check if .Trash dir exists, if yes, use it.  else create
	// .Trash.$uid dir
	legacyTrashPath := filepath.Join(ot.topDir, ".Trash")
	uid := os.Geteuid()
	newTrashDir := fmt.Sprintf(".Trash.%d", uid)

	legacyExists, err := directoryExists(legacyTrashPath)
	if err != nil || !legacyExists {
		if err := os.Mkdir(newTrashDir, 0o700); err != nil {
			return setupTrashDir(newTrashDir)
		} else {
			return "", err
		}
	}
	return "", err
}

func setupTrashDir(path string) (string, error) {
	subDirs := []string{
		filepath.Join(path, "files"),
		filepath.Join(path, "info"),
	}
	for _, dir := range subDirs {
		// 0700 is preferred for private trash folders
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return "", fmt.Errorf("failed to create trash subdir %s: %w", dir, err)
		}
	}
	dirCachefile := filepath.Join(path, "directorysizes")
	if err := makeCacheFile(dirCachefile); err != nil {
		return "", fmt.Errorf("failed to setup cache file: %w", err)
	}
	return path, nil
}

// makeCacheFile creates cache file on its absence, returns any errors
// in the process.
func makeCacheFile(fileAbsPath string) error {
	// os.O_CREATE: Create file if it doesn't exist
	// os.O_EXCL: Return error if file already exists
	// os.O_WRONLY: Open for writing only
	f, err := os.OpenFile(fileAbsPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o666)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil // File already exists, nothing to do
		}
		// Wrapping the error is better for debugging
		return fmt.Errorf("error creating cache file %s: %w", fileAbsPath, err)
	}
	// Always check for errors when closing a file you wrote to
	return f.Close()
}

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

// func trashAFile(file string) string {
// }

// func trashDipatch(file string) string {
// }
