package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/url"
	"time"

	"github.com/debuqer/telem/src/domains/models"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

type Form struct {
	Action string
	Method string
}

type RegisterPageData struct {
	Title string
	Form  Form
}

func (controller *UserController) Register(c *gin.Context) {
	location := url.URL{Path: "/register"}
	data := RegisterPageData{
		Title: "Registeration",
		Form: Form{
			Action: location.RequestURI(),
			Method: "POST",
		},
	}

	tmpl, err := template.ParseFiles("src/domains/templates/auth/register.html")
	if err != nil {
		fmt.Println(err)
	}

	tmpl.Execute(c.Writer, data)
}

func (controller *UserController) DoRegister(c *gin.Context) {
	user := models.User{
		Id:        0,
		Username:  c.Request.FormValue("username"),
		Password:  c.Request.FormValue("password"),
		CreatedAt: time.Now(),
	}

	err := user.Insert()
	if err != nil {
		log.Fatal(err)
	}

	location := url.URL{Path: "/register"}
	c.Redirect(302, location.RequestURI())
}
