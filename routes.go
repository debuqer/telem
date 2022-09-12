package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"telem/helpers"
	"telem/models"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u, _ := models.CurrentUser(r)

	err := tpl.ExecuteTemplate(w, "signup.gohtml", struct {
		Title string
		Csrf  string
		Data  struct {
			User models.User
		}
	}{
		"Post",
		helpers.GetCsrfToken(w, r),
		struct{ User models.User }{
			u,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func applySignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	profileUrl, err := models.UploadProfile(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}

	err = models.AddUser(r.FormValue("name"), r.FormValue("username"), r.FormValue("password"), profileUrl)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u, err := models.CurrentUser(r)
	if err == nil {
		c, err := r.Cookie("session")
		cv := helpers.GetCookieValue(c, err)
		if cv == "" {
			uuid, _ := uuid.NewUUID()
			helpers.SetCookie(w, "session", uuid.String())
			helpers.SetSession(uuid.String(), "username", u.Username)
		} else {
			helpers.SetSession(cv, "username", u.Username)
		}

		http.Redirect(w, r, "/feed", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(w, "login.gohtml", struct {
		Title string
		Csrf  string
	}{
		"Login / Signup",
		helpers.GetCsrfToken(w, r),
	})
}

func applyLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	u, err := models.UserLogin(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !helpers.MatchCsrf(r, r.FormValue("csrf_token")) {
		http.Error(w, "Csrf Not Found", http.StatusUnauthorized)
		return
	}

	sid := uuid.Must(uuid.NewRandom()).String()
	helpers.SetCookie(w, "session", sid)
	helpers.SetSession(sid, "username", u.Username)

	http.Redirect(w, r, "/feed", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c, err := r.Cookie("session")

	if err == nil {
		c.MaxAge = -1
		http.SetCookie(w, c)
		helpers.UnsetSession(c.Value, "sid")
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}

func feed(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u, err := models.CurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts, _ := models.GetPosts(0)
	err = tpl.ExecuteTemplate(w, "feed.gohtml", struct {
		Title string
		Csrf  string
		Data  struct {
			User  models.User
			Posts models.Posts
		}
	}{
		"Feed",
		helpers.GetCsrfToken(w, r),
		struct {
			User  models.User
			Posts models.Posts
		}{
			u,
			posts,
		},
	})
	if err != nil {
		panic(err)
	}
}

func score(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, _ := models.CurrentUser(r)
	pid := r.FormValue("pid")

	pidNumber, _ := strconv.Atoi(pid)
	post := models.FindPost(pidNumber)
	value := r.FormValue("value")
	valueNumber, _ := strconv.Atoi(value)

	post.Score(user, valueNumber)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u, _ := models.CurrentUser(r)
	content := r.FormValue("content")
	pid, _ := strconv.Atoi(r.FormValue("pid"))

	if !helpers.MatchCsrf(r, r.FormValue("csrf_token")) {
		http.Error(w, "Csrf Not Found", http.StatusUnauthorized)
		return
	}

	models.AddPost(u, content, pid)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func singlePost(w http.ResponseWriter, r *http.Request, h httprouter.Params) {
	u, _ := models.CurrentUser(r)
	pid := h.ByName("pid")
	pidNumber, _ := strconv.Atoi(pid)
	post := models.FindPost(pidNumber)

	err := tpl.ExecuteTemplate(w, "single-post.gohtml", struct {
		Title string
		Csrf  string
		Data  struct {
			User models.User
			Post models.Post
		}
	}{
		"Post",
		helpers.GetCsrfToken(w, r),
		struct {
			User models.User
			Post models.Post
		}{
			u,
			post,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}
