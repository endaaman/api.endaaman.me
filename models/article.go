package models

import (
	// "errors"
	// "strconv"
)

type Article struct {
	category Category
	title string
	aliases []string
	tags []string
	image string
	digest string
	priority int
	private bool
	special bool
	date string // 2020-01-01
}

func init() {
}

// func (_ *Article) TableName() string {
//     return "articles"
// }

func (a *Article) FromText(text string, slug string) {
}

func (a *Article) ToText() string {
	return ""
}
