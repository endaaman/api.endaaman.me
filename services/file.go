package services

import (
	"fmt"
	"mime/multipart"

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

func SaveToFile(rel string, file multipart.File) error {
	return infras.SaveToFile(rel, file)
}

func MoveFile(rel, dest string) error {
	return infras.RenameFile(rel, dest)
}
