package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/sys/unix"
)

var (
	mountCache = make(map[string]string)
	cacheMu    sync.RWMutex
)

func GetMountPoint(path string) (string, error) {
	absPath, err := fileCheck(path)
	if err != nil {
		return "", nil
	}
	dir := filepath.Dir(absPath)
	// Quick Read Lock check
	cacheMu.RLock()
	if mp, ok := mountCache[dir]; ok {
		cacheMu.RUnlock()
		return mp, nil
	}
	cacheMu.RUnlock()

	mp, err := getDirectMountPoint(path)
	if err != nil {
		return "", err
	}
	// Write to cache
	cacheMu.Lock()
	mountCache[dir] = mp
	cacheMu.Unlock()

	return mp, nil
}

func getDirectMountPoint(path string) (string, error) {
	absPath, err := fileCheck(path)
	if err != nil {
		return "", err
	}
	// Get device ID for the current path
	var stat unix.Stat_t
	if err := unix.Stat(absPath, &stat); err != nil {
		return "", err
	}
	dev := stat.Dev
	current := absPath
	for {
		parent := filepath.Dir(current)
		// Stop if we hit the root directory
		if current == parent {
			return current, nil
		}
		var pStat unix.Stat_t
		if err := unix.Stat(parent, &pStat); err != nil {
			return "", err
		}
		// If the parent has a different Device ID, 'current' is the
		// mount point
		if pStat.Dev != dev {
			return current, nil
		}
		current = parent
	}
}

// Parsing /proc/mounts
type Mount struct {
	Device  string
	Path    string
	FSType  string
	Options string
}

func (m Mount) String() string {
	// return fmt.Sprintf("%s %s %s %s", m.Device, m.Path, m.FSType, m.Options)
	return fmt.Sprintf("%s %s %s", m.Device, m.Path, m.FSType)
}

var MountFile = "/proc/mounts"

func ParseMountPoints() ([]Mount, error) {
	mfile, err := os.Open(MountFile)
	if err != nil {
		return nil, err
	}
	defer mfile.Close()

	var mounts []Mount
	scanner := bufio.NewScanner(mfile)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		fsType := fields[2]

		switch fsType {
		case "proc", "sysfs", "tmpfs", "devtmpfs":
			continue
		}

		mountPath := fields[1]
		// Dealing the octal escapes (\040)
		if unquoted, err := strconv.Unquote(`"` + mountPath + `"`); err == nil {
			mountPath = unquoted
		}

		mounts = append(mounts, Mount{
			Device:  fields[0],
			Path:    mountPath,
			FSType:  fsType,
			Options: fields[3],
		})
	}
	return mounts, scanner.Err()
}
