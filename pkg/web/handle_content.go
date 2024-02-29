package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-osin/osin"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/models/content"
)

const (
	LimitArticle = 3
	LimitLinks   = 10
)

func (s *server) welcome(c *gin.Context) {

	articles, err := backends.LoadArticles(LimitArticle, 0)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	links, err := backends.LoadLinks(LimitLinks, 0)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	//execute the template
	s.Render(c, "welcome.html", map[string]interface{}{
		"ctx":      c,
		"articles": articles,
		"links":    links,
	})
}

func (s *server) articleView(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if id < 1 {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	article, err := backends.LoadArticle(id)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	fmt.Printf("%s\n%s\n", article.HtmlTitle(), article.HtmlContent())

	if IsAjax(c.Request) {
		res := make(osin.ResponseData)
		res["data"] = article
		c.JSON(http.StatusOK, res)
		return
	}

	s.Render(c, "article_view.html", map[string]interface{}{
		"ctx":     c,
		"article": article,
	})
}

func (s *server) articleForm(c *gin.Context) {
	articles, err := backends.LoadArticles(9, 0)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	s.Render(c, "article_edit.html", map[string]interface{}{
		"ctx":      c,
		"articles": articles,
	})
}

func (s *server) articlePost(c *gin.Context) {
	req := c.Request
	obj := new(content.Article)
	err := binding.FormPost.Bind(req, obj)
	if err != nil {
		log.Printf("bind %v to obj ERR: %s", req.PostForm, err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	user := UserWithContext(c)
	obj.Author = user.UID
	res := make(osin.ResponseData)
	err = backends.SaveArticle(obj)
	if err == nil {
		res["ok"] = true
	} else {
		res["ok"] = false
		log.Printf("save article ERR %s", err)
	}
	c.JSON(http.StatusOK, res)
}

func (s *server) linksForm(c *gin.Context) {
	links, err := backends.LoadLinks(9, 0)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	s.Render(c, "links.html", map[string]interface{}{
		"ctx":   c,
		"links": links,
	})
}

func (s *server) linksPost(c *gin.Context) {

	req := c.Request
	res := make(osin.ResponseData)
	if req.FormValue("op") == "new" {
		obj := new(content.Link)
		err := binding.FormPost.Bind(req, obj)
		if err != nil {
			log.Printf("bind %v to obj ERR: %s", req.PostForm, err)
			res["ok"] = false
			res["error"] = map[string]string{"message": err.Error()}
		}
		user := UserWithContext(c)
		obj.Author = user.UID
		err = backends.SaveLink(obj)
		if err != nil {
			res["ok"] = false
			res["error"] = map[string]string{"message": err.Error()}
		} else {
			res["ok"] = true
		}
	} else {
		pk, name, value := req.PostFormValue("pk"), req.PostFormValue("name"), req.PostFormValue("value")
		// log.Printf("new post: pk %s, name %s, value %s", pk, name, value)
		if pk == "" {
			res["ok"] = false
			res["error"] = map[string]string{"message": "pk is empty"}
			c.JSON(http.StatusOK, res)
			return
		}
		id, err := strconv.Atoi(pk)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		link, err := backends.LoadLink(id)
		if err != nil {
			res["ok"] = false
			res["error"] = map[string]string{"message": "pk is invalid or not found"}
			c.JSON(http.StatusOK, res)
			return
		}
		switch name {
		case "title":
			link.Title = value
		case "url":
			link.SetUrl(value)
		case "position":
			link.Position, _ = strconv.Atoi(value)
		}

		err = backends.SaveLink(link)
		if err != nil {
			res["ok"] = false
			res["error"] = map[string]string{"message": err.Error()}
		} else {
			res["ok"] = true
		}
	}

	c.JSON(http.StatusOK, res)
}
