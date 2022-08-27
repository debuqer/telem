package main

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template
var sqlSrc string

func init() {
	tpl = template.Must(template.New("").ParseGlob("templates/*.gohtml"))
	sqlSrc = "root:@tcp(127.0.0.1:3306)/telem"
}

func main() {
	mux := httprouter.New()

	mux.NotFound = http.HandlerFunc(notFound)
	mux.GET("/signup", signup)
	mux.POST("/signup", applySignup)
	mux.GET("/login", login)
	mux.POST("/login", applyLogin)
	mux.GET("/logout", logout)
	mux.GET("/panel", panel)

	mux.ServeFiles("/statics/*filepath", http.Dir("statics"))
	mux.ServeFiles("/uploads/*filepath", http.Dir("uploads"))
	http.ListenAndServe(":8080", mux)
}
