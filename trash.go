// This files deals with moving and writing trash files
package main

type TrashImpl interface {
	TrashDir() (string, error)
	SourceDir() string
}

func (ht *HomeTrash) SourceDir() string {
	return ht.homeDir
}

func (st *SpecialTrash) SourceDir() string {
	return st.rootDir
}

func Rename(target string, trash TrashImpl) error {
	return nil
}
