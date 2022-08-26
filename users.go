package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type User struct {
	Name       string
	Username   string
	ProfileUrl string
	Password   string
}

var users []User

func addUser(u User) error {
	if u.Name == "" || len(u.Name) < 3 {
		return errors.New("Name must contains at least 3 characters")
	}
	if u.Username == "" || len(u.Username) < 3 {
		return errors.New("Username must contains at least 3 characters")
	}

	if !isUsernameUnique(u.Username) {
		return errors.New("Username must be unique")
	}

	if u.Password == "" || len(u.Password) < 6 {
		return errors.New("Password must contains at least 6 characters")
	}

	users = append(users, u)
	fmt.Println("user", u.Name, "with username", u.Username, "created")

	return nil
}

func getUser(i int) User {
	return users[i]
}

func checkPassword(u User, p string) bool {
	return u.Password == p
}

func isUsernameUnique(username string) bool {
	for _, element := range users {
		if username == element.Username {
			return false
		}
	}

	return true
}

func uploadProfile(r *http.Request) (string, error) {

	f, h, err := r.FormFile("profile")
	if err != nil {
		return "", errors.New("Profile image must be present")
	}

	profileUrl := "uploads/" + h.Filename

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return "", errors.New("Profile image cant be read")
	}

	newFile, err := os.Create("uploads/" + h.Filename)
	if err != nil {
		return "", errors.New("Profile image does not seted properly")
	}
	defer newFile.Close()
	newFile.Write(bs)

	return profileUrl, nil
}
