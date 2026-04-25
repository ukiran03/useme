package fsys

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

type MountInfo struct {
	DevID      uint64 // Device ID 8:12 -> uint64
	MountPoint string // e.g., /mnt/data
	FSType     string // e.g., ext4 (to skip network mounts like nfs/smb)
	IsReadOnly bool   // Parsed from vfsOpts, sbOpts
}

func LoadMountInfo(
	r io.Reader, filterFunc FilterFunc, minfomap map[uint64]*MountInfo,
) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		fields := strings.Fields(text)
		numFields := len(fields)
		if numFields < 10 {
			// should be at least 10 fields
			return fmt.Errorf(
				"parsing '%s' failed: not enough fields (%d)", text, numFields,
			)
		}
		sepIdx := numFields - 4 // "-"
		major, minor, ok := strings.Cut(fields[2], ":")
		if !ok {
			return fmt.Errorf(
				"parsing '%s' failed: unexpected major:minor pair %s",
				text, fields[2],
			)
		}
		vfsOpts := fields[5]
		sbOpts := fields[numFields-1]
		minfo := &MountInfo{
			DevID:      unix.Mkdev(uint32(toInt(major)), uint32(toInt(minor))),
			MountPoint: unescape(fields[4]),
			FSType:     unescape(fields[sepIdx+1]),
			IsReadOnly: slices.Contains(strings.Split(vfsOpts, ","), "ro") ||
				slices.Contains(strings.Split(sbOpts, ","), "ro"),
		}
		if filterFunc != nil {
			if ignorable := filterFunc(minfo); ignorable {
				continue
			}
		}
		minfomap[minfo.DevID] = minfo

	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// toInt converts a string to an int, and ignores any numbers parsing
// errors, as there should not be any.
func toInt(s string) int { i, _ := strconv.Atoi(s); return i }
