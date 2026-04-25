package trash

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/sys/unix"
	"ukiran.com/urm/internal/env"
)

// getTrashPath: returns the valid TrashPath for the given mountpoint(topDir),
// for Home: return XDG_DATA_HOME/Trash
// for Non-Home: M01 or M02
func getTrashPath(topDir string, atHome bool) string {
	if atHome {
		trashHome := os.Getenv("XDG_DATA_HOME")
		if trashHome == "" {
			trashHome = filepath.Join(env.HomeDir, ".local", "share")
		}
		return filepath.Join(trashHome, "Trash")
	}

	// Method (1): $topdir/.Trash/$uid
	// Spec requirement: $topdir/.Trash must be a directory (not a symlink) and
	// must have the sticky bit set.
	dotTrash := filepath.Join(topDir, ".Trash")
	if info, err := os.Lstat(dotTrash); err == nil && info.IsDir() {
		// check sticky bit
		if (info.Mode() & os.ModeSticky) != 0 {
			M01Trash := filepath.Join(dotTrash, strconv.Itoa(env.UID))
			if info, err := os.Lstat(M01Trash); err == nil && info.IsDir() {
				// check write permissions
				if unix.Access(M01Trash, unix.W_OK|unix.X_OK) == nil {
					return M01Trash
				}
			}
		}
	}

	// Method (2): Per-User trash $topdir/.Trash-$uid
	// This is the fallback if Method 1 fails or .Trash doesn't exist.
	M02Trash := filepath.Join(topDir, fmt.Sprintf(".Trash-%d", env.UID))
	return M02Trash
}

// Helper to handle the directory creation logic
func ensureTrashDir(path string) (string, error) {
	if err := os.MkdirAll(path, 0o700); err != nil {
		return "", &TrashError{Op: "mkdir", Path: path, Err: err}
	}
	return path, nil
}
