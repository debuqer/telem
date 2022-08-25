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

	tpl.ExecuteTemplate(w, "login.gohtml", data)
}

func applyLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := struct {
		Username string
		Name     string
	}{
		r.FormValue("username"),
		r.FormValue("name"),
	}

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

	tpl.ExecuteTemplate(w, "login.gohtml", data)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}
