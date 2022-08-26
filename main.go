package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").ParseGlob("templates/*.gohtml"))
}

func main() {
	mux := httprouter.New()

	mux.NotFound = http.HandlerFunc(notFound)
	mux.GET("/signup", signup)
	mux.POST("/signup", applySignup)
	mux.GET("/login", login)
	mux.POST("/login", applyLogin)
	mux.GET("/logout", logout)

	mux.ServeFiles("/statics/*filepath", http.Dir("statics"))
	http.ListenAndServe(":8080", mux)
}

func signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func applySignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	profileUrl, err := uploadProfile(r)
	if err != nil {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		fmt.Println(err)
	}

	err = addUser(User{
		Name:       r.FormValue("name"),
		Username:   r.FormValue("username"),
		Password:   r.FormValue("password"),
		ProfileUrl: profileUrl,
	})

	if err != nil {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		fmt.Println(err)
	}

	http.Redirect(w, r, "/signup", http.StatusSeeOther)
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
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}
