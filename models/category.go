package models

type Category struct {
	Base
	Slug string `json:"slug"`
	Name string `json:"name"`
	Identified bool
}

type CategoryMeta struct {
	Name string `json:"name"`
}

func NewCategory() *Category {
	c := Category{}
	return &c
}
