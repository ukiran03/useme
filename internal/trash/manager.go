package trash

// TODO: obsolete type
type TrashManager struct {
	HomeTrash *TrashCan
	Mounts    map[string]*TrashCan // Keyed by mount point path
}

type MoveStrategy int

const (
	MoveAtomic   MoveStrategy = iota // Standard os.Rename
	MoveFallback                     // Copy to Home + Delete source
	MoveIdentify                     // Just identify, don't move (dry run)
)

// TODO: obsolete func
func (tm *TrashManager) Put(filent *TrashEntry, strategy MoveStrategy) error {
	targetCan, err := tm.FindTarget(filent)
	if err != nil {
		switch strategy {
		case MoveFallback:
			// Prompt to user and operate
			return tm.CopyAndDelete(filent)
		}
	}
	return targetCan.Move(filent)
}

func (tm *TrashManager) CopyAndDelete(filent *TrashEntry) error {
	// TODO: Copy to HomeTrash and delete the source file (filent)
	panic("unimplemented")
}
