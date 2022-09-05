package main

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func FormatDate(t time.Time) string {
	return t.Format("02 Jan 2006 15:04")
}

func init() {
	fm := template.FuncMap{
		"FormatDate": FormatDate,
	}

	tpl = template.Must(template.New("").Funcs(fm).ParseGlob("templates/layouts/*.gohtml"))
	tpl.ParseGlob("templates/*.gohtml")
}

func main() {
	port := "8000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	mux := httprouter.New()

	mux.NotFound = http.HandlerFunc(notFound)
	mux.GET("/signup", signup)
	mux.POST("/signup", applySignup)
	mux.GET("/login", login)
	mux.POST("/login", applyLogin)
	mux.GET("/logout", logout)
	mux.GET("/feed", feed)
	mux.POST("/post", post)
	mux.POST("/score", score)

	mux.ServeFiles("/statics/*filepath", http.Dir("statics"))
	mux.ServeFiles("/uploads/*filepath", http.Dir("uploads"))
	err := http.ListenAndServe(":"+port, mux)
	panic(err)
}
