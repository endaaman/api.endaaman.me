package models

import (
	"time"
)

type File struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mtime"`
	IsDir   bool      `json:"isDir"`
}

func NewFile(name string, size int64, mtime time.Time, isDir bool) *File {
	f := File{}
	f.Name = name
	f.Size = size
	f.ModTime = mtime
	f.IsDir = isDir
	return &f
}
