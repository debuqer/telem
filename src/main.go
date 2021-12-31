package main

import (
	"fmt"
	"log"

	"github.com/debuqer/telem/src/controllers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func GetEnvVariables() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	GetEnvVariables()

	store := cookie.NewStore([]byte("sessions"))
	r := gin.Default()

	r.Use(sessions.Sessions("secret", store))
	uc := controllers.UserController{}

	r.GET("/hello", func(c *gin.Context) {
		fmt.Println("hello")
	})
	r.GET("/register", uc.Register)
	r.POST("/register", uc.DoRegister)

	r.Run()
}
