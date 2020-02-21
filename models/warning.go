package models

import (
)

type Warning struct {
	Base
	Path string `json:"slug"`
	Message string `json:"name"`
}

func NewWarning() *Warning {
	w := Warning{}
	return &w
}
