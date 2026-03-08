package trash

import (
	"golang.org/x/sys/unix"
	"ukiran.com/useme/internal/fs"
)

type TrashCan struct {
	RootPath string // Location where the Trash folder would be
	DeviceID uint64 // Device ID of the RootPath on the File System
}

func NewTrashCan(minfo *fs.MountInfo) *TrashCan {
	return &TrashCan{
		RootPath: minfo.MountPoint,
		DeviceID: unix.Mkdev(uint32(minfo.Major), uint32(minfo.Minor)),
	}
}

// All these methods of TrashCan assume the relevant directories are
// created and checked in prior

func (tc *TrashCan) Put(entry *FileEntry) error {
}

func (tc *TrashCan) Restore(entry *FileEntry, dst string) error

func (tc *TrashCan) Delete(entry *FileEntry) error

func (tc *TrashCan) List() ([]*FileEntry, error)
