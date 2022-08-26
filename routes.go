package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func applySignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	profileUrl, err := uploadProfile(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}

	err = addUser(User{
		Name:       r.FormValue("name"),
		Username:   r.FormValue("username"),
		Password:   []byte(r.FormValue("password")),
		ProfileUrl: profileUrl,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := currentUser(r)
	if err == nil {
		http.Redirect(w, r, "/panel", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func applyLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	u, err := userLogin(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sid := uuid.Must(uuid.NewRandom()).String()
	SetCookie(w, "session", sid)
	setSession(sid, u.Username)

	http.Redirect(w, r, "/panel", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c, err := r.Cookie("session")

	if err == nil {
		c.MaxAge = -1
		http.SetCookie(w, c)
		unsetSession(c.Value)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}

func panel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u, err := currentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tpl.ExecuteTemplate(w, "panel.gohtml", u)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}
