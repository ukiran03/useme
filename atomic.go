package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type SystemTrash struct {
	trash TrashImpl // TrashImpl is an interface
}

func NewSystemTrash(impl TrashImpl) *SystemTrash {
	return &SystemTrash{
		trash: impl,
	}
}

// atomicTrashOperation does few things in atomic fashion like a
// transaction:
//   - makes a file.trashinfo inside $trash/info for given file
//   - moves (renames) the target file to $trash/files
//   - update the $trash/directorysizes
func (syst *SystemTrash) AtomicTrashOperation(targetFile string) (err error) {
	var infoFile, dstPath, srcPath string
	var isInfoFileCreated, isTargetFileMoved bool

	defer func() {
		if err != nil {
			// Rollbacks
			if isInfoFileCreated {
				os.Remove(infoFile)
			}
			if isTargetFileMoved {
				os.Rename(dstPath, srcPath)
			}
		}
	}()

	// Action-1: Create info file
	// Action-2: Move target file to trash
	return nil
}

func (syst *SystemTrash) makeInfoFile(absFilePath string) (string, error) {
	trashdir := syst.trash.SourceDir()
	filename := filepath.Base(absFilePath)
	uniqName := syst.getUniqTrashName(trashdir, filename)
	infoPath := filepath.Join(trashdir, "info", uniqName+".trashinfo")
	if err := syst.writeInfoFile(infoPath); err != nil {
		return "", err
	}
	// [04-03-2026] TODO: Start here
	return nil
}

func (syst *SystemTrash) writeInfoFile(infopath string) error {
}

func (syst *SystemTrash) getUniqTrashName(trashdir, filename string) string {
	dst := filepath.Join(trashdir, "files", filename)
	// If it doesn't exist, we are good to go
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		return filename
	}
	// If exists, start incrementing
	for i := 2; ; i++ {
		newName := fmt.Sprintf("%s_%d", filename, i) // file.txt_2
		if _, err := os.Stat(newName); os.IsNotExist(err) {
			return newName
		}
	}
}

// https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/4/html-single/introduction_to_system_administration/index#s2-storage-fs-mounting
