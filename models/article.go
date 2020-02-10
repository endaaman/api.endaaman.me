package models

import (
	"fmt"
	// "errors"
	// "encoding/json"
	"net/url"
	"time"
	"strings"
	// "regexp"
	// "strconv"
	"github.com/goccy/go-yaml"

	"github.com/thedevsaddam/govalidator"
)


const HEADER_DELIMITTER = "---"

type Header struct {
	Title string      `json:"title"`
	Tags []string     `json:"tags"`
	Aliases []string  `json:"aliases"`
	Digest string     `json:"digest"`
	Image string      `json:"image"`
	Private bool      `json:"private"`
	Special bool      `json:"special"`
	Priority int      `json:"priority"`
	Date string       `json:"date"`
}


type Article struct {
	Header
	Category string `json:"category"`
	Body string     `json:"body"`
	Warning string  `json:"warning"`
}

func init() {
	govalidator.AddCustomRule("strict_date_str", func(field string, rule string, message string, value interface{}) error {
		s := value.(string)
		fmt.Println("DATE: ", s)
		if s == "" {
			return nil
		}
		_, err := time.Parse("2006-01-02", s)
		if err != nil {
			return fmt.Errorf("`%s` is not valid date", s)
		}
		return nil
	})
}

func (a *Article) FromText(text string, slug string, date string) {
	fmt.Printf("FROM: slug: %s %s\n", slug, date)

	splitted_slug := strings.SplitN(slug, "/", 2)
	if len(splitted_slug) == 2 {
		a.Title = splitted_slug[0]
		a.Category = splitted_slug[1]
	} else {
		a.Title = slug
		a.Category = ""
	}
	a.Date = date

	// var header []string
	lines := strings.Split(text, "\n")
	hasHeaderStart := lines[0] == HEADER_DELIMITTER
	headerEndingLine := -1
	if (hasHeaderStart) {
		for i, line := range lines[1:] {
			if line == HEADER_DELIMITTER {
				headerEndingLine = i + 1
			}
		}
	}
	if !(hasHeaderStart && headerEndingLine > 0) {
		a.Body = text
		return
	}

	a.Body = strings.Join(lines[headerEndingLine+1:len(lines)], "\n")
	headerText := strings.Join(lines[1:headerEndingLine], "\n")
	err := yaml.Unmarshal([]byte(headerText), &a.Header)
	if err != nil {
		a.Warning = err.Error()
	}
}

func (a *Article) Validate() url.Values {
	rules := govalidator.MapData{
		"title": []string{"required"},
		"date": []string{"strict_date_str"},
	}

	opts := govalidator.Options{
		Data:  a,
		Rules: rules,
	}

	v := govalidator.New(opts)

	e := v.ValidateStruct()
	if len(e) > 0 {
		return e
	}
	return nil
}

func (a *Article) ToText() string {
	return ""
}
