package models

import (
	// "errors"
	// "strconv"
	// "time"
)

var articles = []*Article{}

type Article struct {
}

func init() {
}

func (_ *Article) TableName() string {
    return "articles"
}

func GetArticles() []*Article {
	return articles
}

func SetArticles() []*Article {
	return articles
}
