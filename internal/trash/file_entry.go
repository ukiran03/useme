package trash

import (
	"io/fs"
	"os"
	"time"
)

// FileEntry is the representation of a file in trash
type FileEntry struct {
	OrigPath     string      // Absolute path of file "prior" trashing
	TrashPath    string      // Absolute path of file "after" trashing
	Name         string      // Original base name of file
	DeletionDate time.Time   // Time of deletion
	MountRoot    string      // Root path of the Mount Point
	Size         int64       // Size of file in bytes
	IsDir        bool        // Indicates if this is a directory
	FileMode     fs.FileMode // Mode of the file
}

func (f *FileEntry) Exists() bool {
	_, err := os.Stat(f.TrashPath)
	return err == nil
}
