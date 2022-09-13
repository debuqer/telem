package models

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"telem/helpers"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id         int
	Name       string
	Username   string
	ProfileUrl string
	Password   []byte
	Role       string
}

func AddUser(name string, username string, password string, profileUrl string) error {
	validator := validator.New()
	err := validator.Var(name, "required,min=3")
	if err != nil {
		return errors.New("Name must contains at least 3 characters")
	}

	err = validator.Var(username, "required,min=3")
	if err != nil {
		return errors.New("Username must contains at least 3 characters")
	}

	err = validator.Var(username, "required,min=6")
	if len(password) < 6 {
		return errors.New("Password must contains at least 6 characters")
	}

	if !IsUsernameUnique(username) {
		return errors.New("Username must be unique")
	}

	pwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	Conn, err := helpers.GetConn()
	if err != nil {
		log.Fatalln(err)
	}
	defer Conn.Close()

	stmt, err := Conn.Prepare("INSERT INTO users (name, username, password, profile_url, role, created_at ) VALUES (?, ?, ?, ?, ?, NOW())")
	_, err = stmt.Exec(name, username, pwd, profileUrl, "admin")
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func UserLogin(username string, password string) (User, error) {
	element, _ := FindUser(username)
	err := bcrypt.CompareHashAndPassword(element.Password, []byte(password))

	if err == nil && element.Username == username {
		return element, nil
	}

	return User{}, errors.New("User not found")
}

func FindUser(username string) (u User, err error) {
	u = User{}
	Conn, err := helpers.GetConn()
	if err != nil {
		log.Fatalln(err)
	}
	defer Conn.Close()
	stmt, err := Conn.Prepare("SELECT id, username, name, profile_url, password FROM users WHERE username = ?")

	res, err := stmt.Query(username)
	if err != nil {
		return u, err
	}
	defer res.Close()
	res.Next()
	res.Scan(&u.Id, &u.Username, &u.Name, &u.ProfileUrl, &u.Password)

	return u, err
}

func IsUsernameUnique(username string) bool {
	Conn, err := helpers.GetConn()
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

func UploadProfile(r *http.Request) (string, error) {

	f, h, err := r.FormFile("profile")
	if err != nil {
		return "", errors.New("Profile image must be present")
	}

	profileUrl := filepath.Join("/uploads/user/", strings.Replace(h.Filename, " ", "-", -1))

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return "", errors.New("Profile image cant be read")
	}

	newFile, err := os.Create("uploads/user/" + h.Filename)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("Profile image does not seted properly")
	}
	defer newFile.Close()
	newFile.Write(bs)

	return profileUrl, nil
}

func CurrentUser(r *http.Request) (User, error) {
	sid := helpers.GetCookieValue(r.Cookie("session"))
	if sid != "" {
		un, hasItem := helpers.GetSession(sid, "username")
		if !hasItem {
			return User{}, errors.New("Not seted session")
		}

		u, err := FindUser(un)
		if err != nil {
			return User{}, errors.New("Not found user")
		}

		return u, nil
	}

	return User{}, errors.New("Not seted sid")
}

func UserLoggined(r *http.Request) bool {
	if _, err := CurrentUser(r); err != nil {
		return false
	}

	return true
}

func havePerm(u User, roleName string) bool {
	return u.Role == roleName
}
