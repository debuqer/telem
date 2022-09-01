package main

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").ParseGlob("templates/layouts/*.gohtml"))
	tpl.ParseGlob("templates/*.gohtml")
}

func main() {
	mux := httprouter.New()

	mux.NotFound = http.HandlerFunc(notFound)
	mux.GET("/signup", signup)
	mux.POST("/signup", applySignup)
	mux.GET("/login", login)
	mux.POST("/login", applyLogin)
	mux.GET("/logout", logout)
	mux.GET("/feed", feed)
	mux.POST("/post", post)

	mux.ServeFiles("/statics/*filepath", http.Dir("statics"))
	mux.ServeFiles("/uploads/*filepath", http.Dir("uploads"))
	err := http.ListenAndServe(":8088", mux)
	panic(err)
}
