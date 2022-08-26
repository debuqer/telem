package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").ParseGlob("templates/*.gohtml"))
}

func main() {
	mux := httprouter.New()

	mux.NotFound = http.HandlerFunc(notFound)
	mux.GET("/login", login)
	mux.POST("/login", applyLogin)
	mux.GET("/logout", logout)

	mux.ServeFiles("/statics/*filepath", http.Dir("statics"))
	http.ListenAndServe(":8080", mux)
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := struct {
		Username string
		Name     string
	}{
		"",
		"",
	}

	cookie, err := r.Cookie("username")
	if err == nil {
		data.Username = cookie.Value
	}

	tpl.ExecuteTemplate(w, "login.gohtml", data)
}

func applyLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	f, h, err := r.FormFile("profile")
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(w, "File is not readable", http.StatusBadGateway)
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(filepath.Join("uploads/", h.Filename), bs, 0644)
	if err != nil {
		http.Error(w, "bad uploaded", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: r.FormValue("username"),
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c, err := r.Cookie("username")

	if err == nil {
		c.MaxAge = -1
		http.SetCookie(w, c)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}
