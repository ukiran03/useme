package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Error: No input file")
	}

	for _, arg := range args {
		mount, err := getFsMountPoint(arg)
		if err != nil {
			log.Printf("%s: %v\n", arg, err)
			continue
		}
		fmt.Printf("%s: %s\n", arg, mount)
	}
}

func getFsMountPoint(file string) (string, error) {
	path, err := fileCheck(file)
	if err != nil {
		return "", err
	}
	// Statfs gets the filesystem type magic number
	var stat unix.Statfs_t
	err = unix.Statfs(path, &stat)
	if err != nil {
		return "", fmt.Errorf("Error: %v\n", err)
	}
	mountPoint := findMountPoint(path)
	fmt.Printf("File:        %s\n", path)
	fmt.Printf("FS Type Hex: 0x%x\n", stat.Type)
	fmt.Printf("Mount Point: %s\n", mountPoint)
	return mountPoint, nil
}

// findMountPoint identifies the root of a filesystem by traversing up
// the directory tree until the Device ID (stat.Dev) changes. In Unix,
// a change in Device ID signifies that we have crossed a mount
// boundary.
func findMountPoint(path string) string {
	dev := getDev(path)
	for {
		parent := filepath.Dir(path)
		// If the parent directory is on a different device,
		// the current path is the mount point.
		if getDev(parent) != dev {
			return path
		}
		// Stop if we've reached the system root (e.g., "/")
		if parent == path {
			return path
		}
		path = parent
	}
}

func getDev(path string) uint64 {
	var stat unix.Stat_t
	unix.Stat(path, &stat)

	return stat.Dev
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
