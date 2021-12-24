package controllers

import (
	"fmt"
	"html/template"
	"os"
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
	data := RegisterPageData{
		Title: "Registeration",
		Form: Form{
			Action: os.Getenv("URL") + "/register",
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

	user.Insert()
}
