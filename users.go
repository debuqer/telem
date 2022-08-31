package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name       string
	Username   string
	ProfileUrl string
	Password   []byte
	Role       string
}

var users []User

func addUser(u User) error {
	validator := validator.New()
	err := validator.Var(u.Name, "required,min=3")
	if err != nil {
		return errors.New("Name must contains at least 3 characters")
	}

	err = validator.Var(u.Username, "required,min=3")
	if err != nil {
		return errors.New("Username must contains at least 3 characters")
	}

	err = validator.Var(u.Username, "required,min=6")
	if len(u.Password) < 6 {
		return errors.New("Password must contains at least 6 characters")
	}

	if !isUsernameUnique(u.Username) {
		return errors.New("Username must be unique")
	}

	u.Password, _ = bcrypt.GenerateFromPassword(u.Password, bcrypt.MinCost)

	Conn, err := getConn()
	if err != nil {
		log.Fatalln(err)
	}
	defer Conn.Close()

	stmt, err := Conn.Prepare("INSERT INTO users (name, username, password, profile_url, role, created_at ) VALUES (?, ?, ?, ?, ?, NOW())")
	_, err = stmt.Exec(u.Name, u.Username, u.Password, u.ProfileUrl, u.Role)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func userLogin(username string, password string) (User, error) {
	element, _ := findUser(username)
	err := bcrypt.CompareHashAndPassword(element.Password, []byte(password))

	if err == nil && element.Username == username {
		return element, nil
	}

	return User{}, errors.New("User not found")
}

func findUser(username string) (u User, err error) {
	u = User{}
	Conn, err := getConn()
	if err != nil {
		log.Fatalln(err)
	}
	defer Conn.Close()
	stmt, err := Conn.Prepare("SELECT username, name, profile_url, password FROM users WHERE username = ?")
	res, err := stmt.Query(username)
	if err != nil {
		return u, err
	}
	defer res.Close()
	res.Next()
	res.Scan(&u.Username, &u.Name, &u.ProfileUrl, &u.Password)

	return u, err
}

func isUsernameUnique(username string) bool {
	Conn, err := getConn()
	if err != nil {
		log.Fatalln(err)
	}
	defer Conn.Close()
	stmt, err := Conn.Prepare("SELECT COUNT(*) as count FROM users WHERE username = ?")
	if err != nil {
		return false
	}
	res, err := stmt.Query(username)

	var count int
	res.Next()
	res.Scan(&count)

	return count == 0
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

func havePerm(u User, roleName string) bool {
	return u.Role == roleName
}
