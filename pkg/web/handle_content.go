package web

import (
	"fmt"
	"log"
	"strconv"

	"github.com/RangelReale/osin"
	"github.com/gin-gonic/gin/binding"

	"lcgc/platform/staffio/pkg/backends"
	"lcgc/platform/staffio/pkg/models"
	. "lcgc/platform/staffio/pkg/settings"
)

const (
	LimitArticle = 3
	LimitLinks   = 10
)

func welcome(ctx *Context) (err error) {

	if Settings.Debug {
		log.Printf("session Name: %s, Values: %d", ctx.Session.Name(), len(ctx.Session.Values))
		log.Printf("ctx User %v", ctx.User)
	}
	articles, err := backends.LoadArticles(LimitArticle, 0)
	if err != nil {
		return err
	}
	links, err := backends.LoadLinks(LimitLinks, 0)
	if err != nil {
		return err
	}

	//execute the template
	return ctx.Render("welcome.html", map[string]interface{}{
		"ctx":      ctx,
		"articles": articles,
		"links":    links,
	})
}

func articleView(ctx *Context) error {
	id, err := strconv.Atoi(ctx.Vars["id"])
	if err != nil {
		return err
	}
	if id < 1 {
		return fmt.Errorf("invalid id %d", id)
	}
	article, err := backends.LoadArticle(id)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n%s\n", article.HtmlTitle(), article.HtmlContent())

	if ctx.IsAjax() {
		res := make(osin.ResponseData)
		res["data"] = article
		return outputJson(res, ctx.Writer)
	}

	return ctx.Render("article_view.html", map[string]interface{}{
		"ctx":     ctx,
		"article": article,
	})
}

func articleForm(ctx *Context) error {
	if !ctx.checkLogin() {
		return nil
	}
	if !ctx.User.IsKeeper() {
		ctx.toLogin()
		return nil
	}
	articles, err := backends.LoadArticles(9, 0)
	if err != nil {
		return err
	}

	return ctx.Render("article_edit.html", map[string]interface{}{
		"ctx":      ctx,
		"articles": articles,
	})
}

func articlePost(ctx *Context) error {
	if !ctx.checkLogin() {
		return nil
	}
	if !ctx.User.IsKeeper() {
		ctx.toLogin()
		return nil
	}
	req := ctx.Request
	obj := new(models.Article)
	err := binding.FormPost.Bind(req, obj)
	if err != nil {
		log.Printf("bind %v to obj ERR: %s", req.PostForm, err)
		return err
	}
	obj.Author = ctx.User.Uid
	res := make(osin.ResponseData)
	err = backends.SaveArticle(obj)
	if err == nil {
		res["ok"] = true
	} else {
		res["ok"] = false
		log.Printf("save article ERR %s", err)
	}
	return outputJson(res, ctx.Writer)
}

func linksForm(ctx *Context) error {
	if !ctx.checkLogin() {
		return nil
	}
	if !ctx.User.IsKeeper() {
		ctx.toLogin()
		return nil
	}
	links, err := backends.LoadLinks(9, 0)
	if err != nil {
		return err
	}

	return ctx.Render("links.html", map[string]interface{}{
		"ctx":   ctx,
		"links": links,
	})
}

func linksPost(ctx *Context) error {
	if !ctx.checkLogin() {
		return nil
	}
	if !ctx.User.IsKeeper() {
		ctx.toLogin()
		return nil
	}

	req := ctx.Request
	res := make(osin.ResponseData)
	if req.FormValue("op") == "new" {
		obj := new(models.Link)
		err := binding.FormPost.Bind(req, obj)
		if err != nil {
			log.Printf("bind %v to obj ERR: %s", req.PostForm, err)
			res["ok"] = false
			res["error"] = map[string]string{"message": err.Error()}
		}
		obj.Author = ctx.User.Uid
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
			return outputJson(res, ctx.Writer)
		}
		id, err := strconv.Atoi(pk)
		if err != nil {
			return err
		}
		link, err := backends.LoadLink(id)
		if err != nil {
			res["ok"] = false
			res["error"] = map[string]string{"message": "pk is invalid or not found"}
			return outputJson(res, ctx.Writer)
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

	return outputJson(res, ctx.Writer)
}
