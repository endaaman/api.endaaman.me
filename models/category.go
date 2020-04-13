package models

import (
	"encoding/json"
	"fmt"
)

type CategoryMeta struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

type Category struct {
	Base
	CategoryMeta
	Slug string `json:"slug"`
}

func NewCategory(slug string) *Category {
	c := Category{}
	c.Slug = slug
	return &c
}

func (c *Category) FromJSON(jsonStr string) error {
	err := json.Unmarshal([]byte(jsonStr), &c.CategoryMeta)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal json: %s", err.Error())
	}
	return nil
}
