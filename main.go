package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").ParseGlob("statics/*.gohtml"))
}

func main() {
	http.HandleFunc("/login", login)

	fmt.Println("Start listening on 8080")
	http.ListenAndServe(":8080", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	type LoginForm struct {
		Username string
	}

	var data LoginForm
	if r.Method == http.MethodPost {
		data = LoginForm{
			r.FormValue("username"),
		}
	} else {
		data = LoginForm{
			"",
		}
	}

	tpl.ExecuteTemplate(w, "login.gohtml", data)
}
