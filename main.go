package main

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").ParseGlob("templates/*.gohtml"))
}

func main() {
	mux := httprouter.New()

	mux.GET("/login-image", loginImage)
	mux.GET("/login", login)
	mux.POST("/login", applyLogin)
	http.ListenAndServe(":8080", mux)
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := struct {
		Username string
	}{
		"",
	}

	tpl.ExecuteTemplate(w, "login.gohtml", data)
}

func loginImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "statics/img/login.jpg")
}

func applyLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := struct {
		Username string
	}{
		r.FormValue("username"),
	}

	tpl.ExecuteTemplate(w, "login.gohtml", data)
}
