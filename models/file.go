package models

import (
	"time"
)

type File struct {
	Name string
	Size int64
	ModTime time.Time
	IsDir bool
}

func NewFile(name string, size int64, mtime time.Time, isDir bool) *File {
	f := File{}
	f.Name = name
	f.Size = size
	f.ModTime = mtime
	f.IsDir = isDir
	return &f
}
