package main

import (
	"fmt"
)

// Unable to find or create trash directory /home/ukiran/Mobile/.Trash-1000
// to trash /home/ukiran/Mobile/InternalSharedStorage/file.txt
func internalError(dst, file string) error {
	return fmt.Errorf(
		"Unable to find or create trash directory %s to trash %s", dst, file,
	)
}
