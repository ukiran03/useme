package trash

type TrashManager struct {
	HomeTrash *TrashCan
	Mounts    map[string]*TrashCan // Keyed by mount point path
}

// FindTarget decides where the trash should go.
func (tm *TrashManager) FindTarget(filePath string) *TrashCan {
	// Get device ID of filePath
	// If same as Home, return tm.HomeTrash
	// Else, lookfor/initialize a TrashCan on that mount point

	return &TrashCan{}
}
