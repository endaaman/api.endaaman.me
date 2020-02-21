package models

import (
	"encoding/json"
)

type Category struct {
	Base
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type CategoryMeta struct {
	Name string `json:"name"`
}

func NewCategory() *Category {
	c := Category{}
	return &c
}

func (c *Category) FromJSON(jsonStr string) error {
	meta := &CategoryMeta{}
	err := json.Unmarshal([]byte(jsonStr), &meta)
	if err != nil {
		return err
	}
	return nil
}
