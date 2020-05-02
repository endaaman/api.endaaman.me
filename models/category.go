package models

import (
	"strings"
)

type Category struct {
	Base
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

func NewCategory(slug string) *Category {
	c := Category{}
	c.Slug = slug
	c.Name = slug
	c.Priority = 0
	return &c
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
