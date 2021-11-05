package i18n

import (
	_ "golang.org/x/text/message/catalog" // for i18n
)

//go:generate gotext -srclang=en update -out=catalog_gen.go -lang=en,zh-hans,zh-hant
