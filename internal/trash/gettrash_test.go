package trash

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetAdminTrashDir(t *testing.T) {
	// 1. Setup a temp directory to act as our "partition root"
	tempBase, err := os.MkdirTemp("", "trash_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempBase)

	trashPath := filepath.Join(tempBase, ".Trash")

	// Test Case 1: Fails if directory doesn't exist
	_, err = getAdminTrashDir(trashPath)
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}

	// Create the directory
	os.Mkdir(trashPath, 0o777)

	// Test Case 2: Fails if Sticky Bit is missing
	_, err = getAdminTrashDir(trashPath)
	if err == nil || !strings.Contains(err.Error(), "sticky bit") {
		t.Errorf("Expected sticky bit error, got: %v", err)
	}

	// Set Sticky Bit (01000 | 0777)
	os.Chmod(trashPath, os.ModeSticky|0o777)

	// Test Case 3: Ownership Check
	// Note: In a local unit test, the folder will be owned by YOU, not root.
	// This will trigger your 'stat.Uid != 0' check.
	_, err = getAdminTrashDir(trashPath)
	if err == nil || !strings.Contains(err.Error(), "not owned by root") {
		t.Errorf("Expected root ownership error, got: %v", err)
	}
}
