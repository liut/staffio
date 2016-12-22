package backends

import (
	"fmt"
	"log"

	"lcgc/platform/staffio/models"
)

func LoadArticle(id int) (*models.Article, error) {
	a := new(models.Article)

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

func LoadArticles(limit, offset int) (data []*models.Article, err error) {
	if limit < 1 {
		limit = 1
	}
	if offset < 0 {
		offset = 0
	}

	str := `SELECT id, title, content, author, created
	   FROM articles ORDER BY created DESC`

	str = fmt.Sprintf("%s LIMIT %d OFFSET %d", str, limit, offset)

	data = make([]*models.Article, 0)
	qs := func(db dber) error {
		return db.Select(&data, str)
	}

	if err := withDbQuery(qs); err != nil {
		return nil, err
	}

	return data, nil
}

func LoadLinks(limit, offset int) (data []*models.Link, err error) {
	if limit < 1 {
		limit = 1
	}
	if offset < 0 {
		offset = 0
	}

	str := `SELECT * FROM links ORDER BY position`

	str = fmt.Sprintf("%s LIMIT %d OFFSET %d", str, limit, offset)

	data = make([]*models.Link, 0)
	qs := func(db dber) error {
		return db.Select(&data, str)
	}

	if err := withDbQuery(qs); err != nil {
		return nil, err
	}

	return data, nil
}

func SaveArticle(a *models.Article) error {
	qs := func(db dbTxer) error {
		log.Printf("save %d", a.Id)
		if a.Id > 0 {
			str := `UPDATE articles SET title = $1, content = $2, updated = CURRENT_TIMESTAMP WHERE id = $3`
			_, err := db.Exec(str, a.Title, a.Content, a.Id)
			if err == nil {
				return nil
			}
			log.Printf("UPDATE article ERR %s", err)
			return err
		}
		str := `INSERT INTO articles(title, content, author) VALUES($1, $2, $3)`
		_, err := db.Exec(str, a.Title, a.Content, a.Author)
		if err == nil {
			return nil
		}
		log.Printf("INSERT article ERR %s", err)
		return err
	}
	return withTxQuery(qs)
}

func LoadLink(id int) (*models.Link, error) {
	link := new(models.Link)
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

func SaveLink(l *models.Link) error {
	qs := func(db dbTxer) (err error) {
		if l.Id > 0 {
			_, err = db.Exec("UPDATE links SET title = $1, url = $2, position = $3 WHERE id = $4", l.Title, string(l.Url), l.Position, l.Id)
		} else {
			_, err = db.Exec("INSERT INTO links(title, url, author) VALUES($1, $2, $3)", l.Title, string(l.Url), l.Author)
		}
		return
	}
	return withTxQuery(qs)
}
