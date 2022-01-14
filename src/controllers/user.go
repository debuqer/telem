package controllers

import (
	"fmt"
	"html/template"
	"net/url"

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

type RegisterValidation struct {
	Name     string `form:"name" json:"username" binding:"required,min=3"`
	Email    string `form:"email" json:"email" binding:"required,min=3"`
	Password string `form:"password" json:"password" binding:"required,min=6"`
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

	location := url.URL{Path: "/user/register"}
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

	validationForm := RegisterValidation{}
	if err := c.ShouldBind(&validationForm); err != nil {
		session.AddFlash(err.Error())
	} else {
		user := models.User{
			Id:       0,
			Name:     c.Request.FormValue("name"),
			Email:    c.Request.FormValue("email"),
			Password: c.Request.FormValue("password"),
		}

		_, err := user.Insert()
		if err != nil {
			session.AddFlash(err.Error())
		}

	}
	session.Save()

	location := url.URL{Path: "/user/register"}
	c.Redirect(302, location.RequestURI())
}
