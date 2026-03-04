// Special Trash Implementation
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/sys/unix"
)

type SpecialTrash struct {
	rootDir  string // mount point
	trashDir string
}

func NewSpecialTrash(path string) *SpecialTrash {
	return &SpecialTrash{rootDir: path}
}

func (st *SpecialTrash) TrashDir() (string, error) {
	adminTrashRoot := filepath.Join(st.rootDir, ".Trash")
	uid := strconv.Itoa(os.Getuid())
	var trashPath string
	var err error

	// Attempt Method 1: /.Trash/$uid
	secure, _ := adminTrashExistsAndSecure(adminTrashRoot)
	// Ignoring error here to trigger fallback logic
	if secure {
		trashPath = filepath.Join(adminTrashRoot, uid)
		var td string
		td, err = setupTrashDir(trashPath)
		if err == nil {
			return td, nil
		}
	}
	// Fallback to Method 2: /.Trash-$uid
	// This runs if Method 1 was insecure OR if setupTrashDir failed.
	trashPath = filepath.Join(st.rootDir, fmt.Sprintf(".Trash-%d", os.Getuid()))
	return setupTrashDir(trashPath)
}

func adminTrashExistsAndSecure(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	// Must be a directory, NOT a symlink
	if !info.IsDir() {
		return false, nil
	}
	// Check the sticky bit:
	// os.ModeSticky is 0x20000000; in chmod terms, it's the "1" in "1777"
	if (info.Mode() & os.ModeSticky) == 0 {
		return false, fmt.Errorf(
			"security risk: sticky bit not set on %s", path,
		)
	}
	// Check ownership (Standard requirement: must be owned by root)
	stat, ok := info.Sys().(*unix.Stat_t)
	if !ok {
		return false, fmt.Errorf(
			"could not get raw syscall.Stat_t for %s", path,
		)
	}
	if stat.Uid != 0 {
		return false, fmt.Errorf(
			"security risk: %s is not owned by root", path,
		)
	}
	return true, nil
}
