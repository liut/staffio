package content

import (
	"html/template"
	"regexp"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

type Article struct {
	Id      int       `sql:"id,pk" json:"id" form:"id"`
	Title   string    `sql:"title,notnull" json:"title" form:"title" binding:"required"`
	Content string    `sql:"content,notnull" json:"content" form:"content" binding:"required"`
	Author  string    `sql:"author" json:"author"`
	Created time.Time `sql:"created" json:"created"`
	Updated time.Time `sql:"updated,nullempty" json:"updated,omitempty"`
}

func (a *Article) HtmlTitle() template.HTML {
	return SanitizeTitle(a.Title)
}

func (a *Article) HtmlContent() template.HTML {
	return MarkdownSanitize(a.Content)
}

type Link struct {
	Id       int          `sql:"id,pk" json:"id" form:"id"`
	Title    string       `sql:"title,notnull" json:"title" form:"title" binding:"required"`
	Url      template.URL `sql:"url,unique,notnull" json:"url" form:"url" binding:"required"`
	Position int          `sql:"position" json:"position" form:"position"`
	Author   string       `sql:"author" json:"author"`
	Created  time.Time    `sql:"created" json:"created"`
}

func (a *Article) StyleName() string {
	switch a.Id % 5 {
	case 1:
		return "primary"
	case 2:
		return "success"
	case 3:
		return "danger"
	case 4:
		return "warning"
	default:
		return "info"
	}
}

func (a *Link) SetUrl(href string) {
	a.Url = template.URL(href)
}

func (a *Link) HtmlTitle() template.HTML {
	return SanitizeTitle(a.Title)
}

func SanitizeTitle(s string) template.HTML {
	unsafe := blackfriday.MarkdownBasic([]byte(s))
	p := bluemonday.StrictPolicy()
	p.AllowElements("b", "i", "u", "small", "strike", "strong", "tt",
		"pre", "code", "sub", "sup")
	return template.HTML(p.SanitizeBytes(unsafe))
}

func MarkdownSanitize(s string) template.HTML {
	unsafe := blackfriday.MarkdownCommon([]byte(s))
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	html := p.SanitizeBytes(unsafe)
	return template.HTML(html)
}
