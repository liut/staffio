package i18n

import (
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Printer ...
type Printer = message.Printer

var (
	enUS   = language.AmericanEnglish
	zhHans = language.MustParse("zh-Hans")
	zhHant = language.MustParse("zh-Hant")

	matcher language.Matcher
)

func init() {
	matcher = language.NewMatcher(message.DefaultCatalog.Languages())
}

// GetTag ...
func GetTag(r *http.Request) language.Tag {
	var lang string
	if s := r.FormValue("lang"); s != "" {
		lang = s
	} else if c, err := r.Cookie("lang"); err == nil {
		lang = c.String()
	} else {
		lang = "zh-hans"
	}
	accept := r.Header.Get("Accept-Language")
	tag, _ := language.MatchStrings(matcher, lang, accept)
	// tag := message.MatchLanguage(lang, accept, "zh-Hans")
	return tag
}

// GetPrinter ...
func GetPrinter(r *http.Request) *message.Printer {
	return message.NewPrinter(GetTag(r))
}

// Field ...
type Field string

func (f Field) Field() string {
	return string(f)
}

func (f Field) ErrorReq(r *http.Request) string {
	return GetPrinter(r).Sprintf("Error:Field validation for '%s' failed ", f.Field())
}
