package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type HomeTrash struct {
	homeDir  string
	trashDir string
}

func NewHomeTrash() (*HomeTrash, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home dir: %w", err)
	}
	return &HomeTrash{homeDir: home}, nil
}

func (ht *HomeTrash) TrashDir() (string, error) {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = filepath.Join(ht.homeDir, ".local", "share")
	}
	trashPath := filepath.Join(dataHome, "Trash")
	return setupTrashDir(trashPath)
}
