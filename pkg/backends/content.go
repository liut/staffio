package backends

import (
	"fmt"

	"github.com/liut/staffio/pkg/models/content"
)

func LoadArticle(id int) (*content.Article, error) {
	a := new(content.Article)

	qs := func(db dber) error {
		return db.Get(a, `SELECT id, title, content, author, created
		 FROM articles WHERE id = $1`, id)
	}
	err := withDbQuery(qs)
	if err == nil {
		return a, nil
	}
	return nil, err
}

func LoadArticles(limit, offset int) (data []*content.Article, err error) {
	if limit < 1 {
		limit = 1
	}
	if offset < 0 {
		offset = 0
	}

	str := `SELECT id, title, content, author, created
	   FROM articles ORDER BY created DESC`

	str = fmt.Sprintf("%s LIMIT %d OFFSET %d", str, limit, offset)

	data = make([]*content.Article, 0)
	qs := func(db dber) error {
		return db.Select(&data, str)
	}

	if err := withDbQuery(qs); err != nil {
		return nil, err
	}

	return data, nil
}

func LoadLinks(limit, offset int) (data []*content.Link, err error) {
	if limit < 1 {
		limit = 1
	}
	if offset < 0 {
		offset = 0
	}

	str := `SELECT * FROM links ORDER BY position`

	str = fmt.Sprintf("%s LIMIT %d OFFSET %d", str, limit, offset)

	data = make([]*content.Link, 0)
	qs := func(db dber) error {
		return db.Select(&data, str)
	}

	if err := withDbQuery(qs); err != nil {
		return nil, err
	}

	return data, nil
}

func SaveArticle(a *content.Article) error {
	qs := func(db dbTxer) error {
		logger().Debugw("save", "article.title", a.Title)
		if a.Id > 0 {
			str := `UPDATE articles SET title = $1, content = $2, updated = CURRENT_TIMESTAMP WHERE id = $3`
			_, err := db.Exec(str, a.Title, a.Content, a.Id)
			if err == nil {
				return nil
			}
			logger().Infow("update article fail", "err", err)
			return err
		}
		var id int
		str := `INSERT INTO articles(title, content, author) VALUES($1, $2, $3) RETURNING id`
		err := db.QueryRow(str, a.Title, a.Content, a.Author).Scan(&id)
		if err != nil {
			logger().Infow("INSERT article failed", "err", err)
			return err
		}
		logger().Infow("new article", "id", id)
		a.Id = id
		return nil
	}
	return withTxQuery(qs)
}

func LoadLink(id int) (*content.Link, error) {
	link := new(content.Link)
	qs := func(db dber) (err error) {
		err = db.Get(link, "SELECT * FROM links WHERE id = $1", id)
		return
	}
	err := withDbQuery(qs)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func SaveLink(l *content.Link) error {
	qs := func(db dbTxer) (err error) {
		if l.Id > 0 {
			_, err = db.Exec("UPDATE links SET title = $1, url = $2, position = $3 WHERE id = $4", l.Title, string(l.Url), l.Position, l.Id)
		} else {
			var id int
			err = db.QueryRow("INSERT INTO links(title, url, author) VALUES($1, $2, $3) RETURNING id", l.Title, string(l.Url), l.Author).Scan(&id)
			if err == nil {
				l.Id = id
			}
		}
		return
	}
	return withTxQuery(qs)
}
