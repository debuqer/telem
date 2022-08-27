package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template
var Conn sql.DB

func init() {
	tpl = template.Must(template.New("").ParseGlob("templates/*.gohtml"))
}

func main() {
	mux := httprouter.New()
	Conn, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/telem")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(Conn)

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
