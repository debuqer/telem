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

	mux.NotFound = http.HandlerFunc(notFound)
	mux.GET("/login", login)
	mux.POST("/login", applyLogin)

	mux.ServeFiles("/statics/*filepath", http.Dir("statics"))
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

func applyLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := struct {
		Username string
	}{
		r.FormValue("username"),
	}

	tpl.ExecuteTemplate(w, "login.gohtml", data)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}
