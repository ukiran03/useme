package fsys

import (
	"fmt"
	"strings"

	"ukiran.com/urm/internal/env"
)

type FilterFunc func(*MountInfo) (ignore bool)

func IgnoreFsFunc(minfo *MountInfo) bool {
	// Check Read-Only
	if minfo.IsReadOnly {
		return true
	}

	if ignoreFsTypes[minfo.FSType] {
		return true
	}
	// We want /tmp, but we don't want '/dev/shm' or '/run/user/1000'
	if minfo.FSType == "tmpfs" && minfo.MountPoint != "/tmp" {
		// If it's a RAM disk but NOT the standard /tmp, ignore it.
		// Note: Check your run/media paths; if your Kindle is tmpfs (rare),
		// you might need to adjust this.
		if !strings.HasPrefix(minfo.MountPoint, "/run/media") {
			return true
		}
	}
	// Check the Prefixes
	for _, prefix := range ignorePrefixes {
		if strings.HasPrefix(minfo.MountPoint, prefix) {
			return true
		}
	}
	return false
}

var (
	ignoreFsTypes = map[string]bool{
		// Kernel & System Virtual Filesystems
		"proc":        true, // Process information pseudo-filesystem
		"sysfs":       true, // Kernel objects/subsystems interface
		"devtmpfs":    true, // Device node management
		"configfs":    true, // Userspace-driven kernel object configuration
		"debugfs":     true, // Kernel debugging interface
		"tracefs":     true, // Infrastructure for kernel tracing
		"binfmt_misc": true, // Support for non-native binary formats
		"fusectl":     true, // Control interface for FUSE filesystems
		"pstore":      true, // Persistent storage for kernel oops/logs
		"devpts":      true, // Pseudo-terminal slave devices
		"autofs":      true, // Kernel-based automounter support
		"cgroup":      true, // Control Groups v1
		"cgroup2":     true, // Control Groups v2
		"efivarfs":    true, // Interface to EFI variable storage
		"hugetlbfs":   true, // Huge page memory support
		"mqueue":      true, // POSIX message queues
		"nsfs":        true, // Namespace filesystems
		"rpc_pipefs":  true, // RPC pipe communication for NFS/SunRPC
		"squashfs":    true, // Compressed read-only FS (often used for Snaps)

		// Container & Overlay Filesystems
		"overlay": true, // Docker/Container layered filesystems

		// Network Filesystems
		"nfs":        true, // Network File System (general)
		"nfs4":       true, // Network File System v4
		"cifs":       true, // Common Internet File System (Samba)
		"smb3":       true, // Server Message Block v3
		"davfs":      true, // WebDAV (Remote web folders)
		"fuse.sshfs": true, // Filesystem over SSH

		// Virtual Desktop / FUSE Portals
		"fuse.gvfsd-fuse": true, // GNOME Virtual File System
		"fuse.portal":     true, // Flatpak/Snap document portals
	}

	ignorePrefixes = []string{
		// System Critical Roots
		"/dev",  // Static and dynamic device nodes
		"/proc", // Process and kernel information
		"/sys",  // System hardware and driver information
		"/boot", // Bootloader, kernels, and initrd

		// User Data & Special Exclusions
		// "/home",             // The /home mount point itself (prevents .Trash here)
		"/home/docker-data", // Personal exclusion for docker volume data

		// Service & Container Data
		"/var/lib/docker", // Internal Docker storage/images
		"/var/run",        // Modern Linux runtime state (symlink to /run)
		"/run/lock",       // Coordinated device/file locking
		"/run/initramfs",  // Root FS used during early boot

		// Specific User Runtime Portals
		fmt.Sprintf("/run/user/%d/doc", env.UID),  // Virtual document portals (FUSE)
		fmt.Sprintf("/run/user/%d/gvfs", env.UID), // Virtual filesystem mounts (FUSE)
	}
)
