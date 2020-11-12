package backends

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liut/staffio/pkg/models/content"
)

func TestArticle(t *testing.T) {
	a := &content.Article{
		Title: "the article title",
		Content: `## article subject
		`,
	}

	err := SaveArticle(a)
	assert.NoError(t, err)

	_, err = LoadArticle(a.Id)
	assert.NoError(t, err)

	data, err := LoadArticles(5, 0)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	a.Title = "change title"
	err = SaveArticle(a)
	assert.NoError(t, err)
}

func TestLink(t *testing.T) {
	l := &content.Link{
		Title: "link name",
		Url:   "http://localhost",
	}
	err := SaveLink(l)
	assert.NoError(t, err)

	_, err = LoadLink(l.Id)
	assert.NoError(t, err)

	data, err := LoadLinks(5, 0)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	l.Title = "change name"
	err = SaveLink(l)
	assert.NoError(t, err)
}
