package main

import (
	"fmt"
	"net/http"
	"strconv"
	"telem/helpers"
	"telem/models"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
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
	}{
		"Login / Signup",
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

	posts, _ := models.GetFeed()
	err = tpl.ExecuteTemplate(w, "feed.gohtml", struct {
		Title string
		Data  struct {
			User  models.User
			Posts models.Posts
		}
	}{
		"Feed",
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

	http.Redirect(w, r, "/feed", http.StatusSeeOther)
}

func post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u, _ := models.CurrentUser(r)
	content := r.FormValue("content")

	models.AddPost(u, content)

	http.Redirect(w, r, "/feed", http.StatusSeeOther)
}

func singlePost(w http.ResponseWriter, r *http.Request, h httprouter.Params) {
	pid := h.ByName("pid")
	pidNumber, _ := strconv.Atoi(pid)
	post := models.FindPost(pidNumber)

	tpl.ExecuteTemplate(w, "single-post.gohtml", post)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}
