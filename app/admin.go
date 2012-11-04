package app

import (
	"code.google.com/p/gorilla/sessions"
	"github.com/hoisie/web"
	"github.com/jmoiron/monet/conf"
	"github.com/jmoiron/monet/template"
)

type AdminPanel interface {
	Render() string
}

var Panels = []AdminPanel{}

func AddPanel(p AdminPanel) {
	Panels = append(Panels, p)
}

var adminBase = template.Base{Path: "admin/base.mandira"}
var CookieStore = sessions.NewCookieStore([]byte(conf.Config.SessionSecret))

func AttachAdmin(url string) {
	// auth
	web.Get(url+"login/", login)
	web.Post(url+"login/", login)
	web.Get(url+"logout/", logout)
	// users
	/*  too much unnecessary work? 
	    web.Get(url + "users/", userList)
	    web.Get(url + "users/edit/(.*)", userEdit)
	    web.Get(url + "users/delete/(.*)", userDelete)
	    web.Get(url + "users/add/", userAdd)
	    web.Post(url + "users/add/", userAddPost)
	*/
	web.Get(url, adminIndex)
}

// Use this on all admin views to ensure the request is authenticated

func RequireAuthentication(ctx *web.Context) bool {
	session, _ := CookieStore.Get(ctx.Request, "monet-session")

	if session.Values["authenticated"] != true {
		ctx.Redirect(302, "/admin/login/")
		return true
	}
	return false
}

// views

func login(ctx *web.Context) string {
	if ctx.Params != nil {
		p := ctx.Params
		if ValidateUser(p["username"], p["password"]) {
			session, _ := CookieStore.Get(ctx.Request, "monet-session")
			session.Values["authenticated"] = true
			session.Save(ctx.Request, ctx.ResponseWriter)
			ctx.Redirect(302, "/admin/")
		}
	}
	return adminBase.Render("admin/login.mandira", ctx.Params, M{"login": true})
}

func logout(ctx *web.Context) string {
	session, _ := CookieStore.Get(ctx.Request, "monet-session")
	session.Values["authenticated"] = false
	session.Save(ctx.Request, ctx.ResponseWriter)
	ctx.Redirect(302, "/admin/login/")
	return ""
}

func adminIndex(ctx *web.Context) string {
	if RequireAuthentication(ctx) {
		return ""
	}

	panels := []M{}

	for i, p := range Panels {
		side := "left"
		if i%2 != 0 {
			side = "right"
		}
		panels = append(panels, M{
			"side":    side,
			"right":   side == "right",
			"content": p.Render(),
		})
	}

	return adminBase.Render("admin/index.mandira", M{
		"Panels": panels,
	})
	/*
		var posts []db.Post
		var unpublished []db.Post
		var pages []db.Page
		db.Posts().Latest(dict{"published": 1}).Limit(indexListSize).Iter().All(&posts)
		db.Posts().Latest(dict{"published": 0}).Limit(indexListSize).Iter().All(&unpublished)
		db.Pages().C.Find(nil).Limit(indexListSize).Iter().All(&pages)

		return adminBase.Render("admin/index.mandira", dict{
			"posts": posts, "unpublished": unpublished, "pages": pages})
	*/
}
