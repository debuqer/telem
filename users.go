package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name       string
	Username   string
	ProfileUrl string
	Password   []byte
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

	if len(u.Password) < 6 {
		return errors.New("Password must contains at least 6 characters")
	}

	u.Password, _ = bcrypt.GenerateFromPassword(u.Password, bcrypt.MinCost)
	fmt.Println(u)

	users = append(users, u)
	fmt.Println("user", u.Name, "with username", u.Username, "created")

	return nil
}

func userLogin(username string, password string) (User, error) {
	for _, element := range users {
		err := bcrypt.CompareHashAndPassword(element.Password, []byte(password))

		if err == nil && element.Username == username {
			return element, nil
		}
	}

	return User{}, errors.New("User not found")
}

func getUser(i int) User {
	return users[i]
}

func findUser(username string) (u User, err error) {
	u = User{}
	for _, element := range users {
		if element.Username == username {
			u = element
			return u, nil
		}
	}
	err = errors.New("No user found")

	return u, err
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

	profileUrl := filepath.Join("uploads/", strings.Replace(h.Filename, " ", "-", -1))

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

func currentUser(r *http.Request) (User, error) {
	sid := getCookieValue(r.Cookie("session"))
	if sid != "" {
		un := getSession(sid)
		u, err := findUser(un)
		if err != nil {
			return User{}, errors.New("Not found user")
		}

		return u, nil
	}

	return User{}, errors.New("Not seted sid")
}
