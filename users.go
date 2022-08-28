package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	Conn, err := sql.Open("mysql", sqlSrc)
	if err != nil {
		log.Fatalln(err)
	}
	defer Conn.Close()
	fmt.Println(Conn)
	query := "INSERT INTO users ( name, username, password, profile_url, created_at ) VALUES ( '" + u.Name + "', '" + u.Username + "', '" + string(u.Password) + "', '" + u.ProfileUrl + "', NOW() )"
	fmt.Println(query)
	_, err = Conn.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("user", u.Name, "with username", u.Username, "created")

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

func getUser(i int) User {
	return users[i]
}

func findUser(username string) (u User, err error) {
	u = User{}
	Conn, err := sql.Open("mysql", sqlSrc)
	if err != nil {
		log.Fatalln(err)
	}
	defer Conn.Close()
	res, err := Conn.Query("SELECT username, name, profile_url, password FROM users WHERE username = '" + username + "'")
	if err != nil {
		return u, err
	}
	defer res.Close()
	res.Next()
	res.Scan(&u.Username, &u.Name, &u.ProfileUrl, &u.Password)

	return u, err
}

func isUsernameUnique(username string) bool {
	Conn, err := sql.Open("mysql", sqlSrc)
	if err != nil {
		log.Fatalln(err)
	}
	defer Conn.Close()
	res, err := Conn.Query("SELECT COUNT(*) as count FROM users WHERE username = '" + username + "'")
	if err != nil {
		return false
	}

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
