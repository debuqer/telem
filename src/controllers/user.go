package controllers

import (
	"fmt"
	"html/template"
	"net/url"
	"time"

	"github.com/debuqer/telem/src/domains/models"
	"github.com/debuqer/telem/src/domains/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

type Form struct {
	Action string
	Method string
}

type RegisterPageData struct {
	Title   string
	Form    Form
	Message services.Message
}

func (controller *UserController) Register(c *gin.Context) {
	session := sessions.Default(c)
	message := services.GetMessageContainer(session)
	session.Save()

	location := url.URL{Path: "/register"}
	data := RegisterPageData{
		Title: "Registeration",
		Form: Form{
			Action: location.RequestURI(),
			Method: "POST",
		},
		Message: message,
	}

	tmpl, err := template.ParseFiles("src/domains/templates/auth/register.html")
	if err != nil {
		fmt.Println(err)
	}

	tmpl.Execute(c.Writer, data)
}

func (controller *UserController) DoRegister(c *gin.Context) {
	session := sessions.Default(c)
	user := models.User{
		Id:        0,
		Username:  c.Request.FormValue("username"),
		Password:  c.Request.FormValue("password"),
		CreatedAt: time.Now(),
	}

	_, err := user.Insert()
	if err != nil {
		session.AddFlash(err.Error())
	}
	session.Save()

	location := url.URL{Path: "/register"}
	c.Redirect(302, location.RequestURI())
}
