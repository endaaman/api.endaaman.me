package models

import (
	// "errors"
	// "strconv"
	// "time"
    "github.com/astaxie/beego/orm"
)

type Article struct {
    Id	int	`orm:"pk;unique;auto;column(article_id)"`
}

func init() {
    orm.RegisterModel(new(Article))
}

func GetAllArticles() []*Article {
	return []*Article{}
}
