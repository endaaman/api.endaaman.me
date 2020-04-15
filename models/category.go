package models

import (
	"encoding/json"
	"fmt"
	"strings"
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

func (a *Category) Compare(b *Category) bool {
	// larger priority goes first
	priDiff := b.Priority - a.Priority
	if priDiff != 0 {
		return priDiff > 0
	}
	// smaller slug goes first
	return strings.Compare(a.Slug, b.Slug) > 0
}
