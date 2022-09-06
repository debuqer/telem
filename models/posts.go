package models

import (
	"fmt"
	"telem/helpers"
	"time"
)

type Post struct {
	Id        int
	User      User
	Content   string
	CreatedAt time.Time
	Likes     int
	Dislikes  int
}

type Posts []Post

func GetFeed() ([]Post, error) {
	posts := Posts{}

	conn, err := helpers.GetConn()
	if err != nil {
		return make(Posts, 0), err
	}

	stmt, err := conn.Prepare("SELECT posts.id, posts.content, posts.created_at, users.name, users.username, users.profile_url, (Select count(1) from likes where post_id = posts.id AND value = 1) as likes, (Select count(1) from likes where post_id = posts.id AND value = -1) as dislikes FROM posts JOIN users ON users.ID = posts.user_id ORDER BY ID DESC LIMIT 10")
	row, err := stmt.Query()

	for row.Next() {
		u := User{}
		p := Post{}
		var cr string

		row.Scan(&p.Id, &p.Content, &cr, &u.Name, &u.Username, &u.ProfileUrl, &p.Likes, &p.Dislikes)
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", cr)

		p.User = u

		posts = append(posts, p)
	}

	return posts, nil
}

func FindPost(pid int) Post {
	conn, err := helpers.GetConn()
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	stmt, err := conn.Prepare("SELECT posts.id, posts.content, posts.created_at, users.name, users.username, users.profile_url, (Select count(1) from likes where post_id = posts.id AND value = 1) as likes, (Select count(1) from likes where post_id = posts.id AND value = -1) as dislikes FROM posts JOIN users ON users.ID = posts.user_id WHERE posts.id = ? ORDER BY ID DESC")
	row, _ := stmt.Query(pid)

	u := User{}
	p := Post{}
	for row.Next() {
		var cr string

		row.Scan(&p.Id, &p.Content, &cr, &u.Name, &u.Username, &u.ProfileUrl, &p.Likes, &p.Dislikes)
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", cr)

		p.User = u
	}

	return p
}

func AddPost(u User, content string) {
	conn, err := helpers.GetConn()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	c := content

	stmt, err := conn.Prepare("INSERT INTO posts ( content, user_id, created_at ) VALUES ( ?, ?, NOW() )")
	stmt.Exec(c, u.Id)
}

func (p *Post) Score(u User, value int) {
	conn, err := helpers.GetConn()
	if err != nil {
		panic(err)
	}
	stmt, err := conn.Prepare("DELETE FROM likes WHERE post_id = ? AND user_id = ?")
	if err != nil {
		panic(err)
	}
	stmt.Exec(p.Id, u.Id)

	stmt, err = conn.Prepare("INSERT INTO likes( user_id, post_id, `value` ) VALUES (?, ?, ?)")
	if err != nil {
		panic(err)
	}
	stmt.Exec(u.Id, p.Id, value)
}
