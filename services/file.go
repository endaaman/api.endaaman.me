package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/endaaman/api.endaaman.me/infras"
	"github.com/endaaman/api.endaaman.me/models"
)

func FileExists(rel string) bool {
	stat := infras.GetStat(rel)
	return stat != nil
}

func IsDir(rel string) bool {
	stat := infras.GetStat(rel)
	return stat != nil && stat.IsDir()
}

func ListDir(rel string) ([]*models.File, error) {
	if !IsDir(rel) {
		return nil, fmt.Errorf("Can not read the path: `%s`", rel)
	}
	return infras.ListDir(rel), nil
}

func DeleteFile(rel string) error {
	return infras.DeleteFile(rel)
}

func SaveFiles(rel string, headers []*multipart.FileHeader) error {
	if len(headers) == 0 {
		return errors.New("Uploaded file is empty")
	}

	m := make(map[string]bool)
	for _, header := range headers {
		name := header.Filename
		if m[name] {
			return fmt.Errorf("Duplicated files(%s) are uploaded", name)
		}
		if !m[name] {
			m[name] = true
		}
		target := filepath.Join(rel, header.Filename)
		if FileExists(target) {
			return fmt.Errorf("The file(%s) already exists.", target)
		}
	}

	for _, header := range headers {
		file, err := header.Open()
		if err != nil {
			return fmt.Errorf("Failed to open file `%s`:  %v", header.Filename, err)
		}
		err = infras.SaveToFile(file, filepath.Join(rel, header.Filename))
		if err != nil {
			return err
		}
	}
	return nil
}

func Mkdir(rel string) error {
	return infras.Mkdir(rel)
}

func MoveFile(rel, dest string) error {
	return infras.RenameFile(rel, dest)
}
