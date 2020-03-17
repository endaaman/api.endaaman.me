package services

import (
	"fmt"
	"mime/multipart"
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/infras"
)

func ListDir(rel string) ([]*models.File, error) {
	if !infras.IsDir(rel) {
		return nil, fmt.Errorf("Can not read the path: `%s`", rel)
	}
	return infras.ListDir(rel), nil
}

func IsDir(rel string) bool {
	return infras.IsDir(rel)
}

func DeleteFile(rel string) error {
	return infras.DeleteFile(rel)
}

func SaveToFile(rel string, file multipart.File) error {
	return infras.SaveToFile(rel, file)
}

func Exists(rel string) bool {
	return infras.Exists(rel)
}

