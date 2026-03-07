package trash

type TrashCan struct {
	RootPath string
	DeviceID uint64
}

func (tc *TrashCan) Put(entry *FileEntry) error

func (tc *TrashCan) Restore(entry *FileEntry, dst string) error

func (tc *TrashCan) Delete(entry *FileEntry) error

func (tc *TrashCan) List() ([]*FileEntry, error)
