package trash

import (
	"ukiran.com/urm/internal/env"
	"ukiran.com/urm/internal/fsys"
)

type MountRegistry struct {
	mounts map[uint64]*fsys.MountInfo
}

func NewMountRegistry() *MountRegistry {
	minfomap := make(map[uint64]*fsys.MountInfo, 64)
	return &MountRegistry{
		mounts: minfomap,
	}
}

func (mr *MountRegistry) GetMountPoint(devID uint64) string {
	if devID == env.HomeDevID {
		return env.HomeDir
	}
	if mountPath, ok := mr.mounts[devID]; ok {
		return mountPath.MountPoint
	}
	// Only parse /proc/self/mountinfo if devID isn't found
	mr.Load()
	return mr.mounts[devID].MountPoint
}

func (mr *MountRegistry) Load() error {
	f, err := env.OpenMountInfo()
	if err != nil {
		return err
	}
	minfomap := make(map[uint64]*fsys.MountInfo, 64)
	err = fsys.LoadMountInfo(f, fsys.IgnoreFsFunc, minfomap)
	if err != nil {
		return err
	}
	return nil
}
