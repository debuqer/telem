package models

import (
	"database/sql"
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
	Replies   int
	Posts     []Post
}

type Posts []Post

func GetPosts(pid int, uid int) ([]Post, error) {
	posts := Posts{}

	conn, err := helpers.GetConn()
	if err != nil {
		return make(Posts, 0), err
	}

	var row *sql.Rows
	var stmt *sql.Stmt

	if pid == 0 {
		if uid != 0 {
			stmt, err = conn.Prepare("SELECT posts.id, posts.content, posts.created_at, users.name, users.username, users.profile_url, (Select count(1) from likes where post_id = posts.id AND value = 1) as likes, (Select count(1) from likes where post_id = posts.id AND value = -1) as dislikes, (select count(1) from posts as p where p.post_id = posts.id) as replies FROM posts JOIN users ON users.ID = posts.user_id WHERE post_id IS NULL AND user_id = ? ORDER BY ID DESC LIMIT 10")
			row, err = stmt.Query(uid)
		} else {
			stmt, err = conn.Prepare("SELECT posts.id, posts.content, posts.created_at, users.name, users.username, users.profile_url, (Select count(1) from likes where post_id = posts.id AND value = 1) as likes, (Select count(1) from likes where post_id = posts.id AND value = -1) as dislikes, (select count(1) from posts as p where p.post_id = posts.id) as replies FROM posts JOIN users ON users.ID = posts.user_id WHERE post_id IS NULL ORDER BY ID DESC LIMIT 10")
			row, err = stmt.Query()
		}
	} else {
		stmt, err = conn.Prepare("SELECT posts.id, posts.content, posts.created_at, users.name, users.username, users.profile_url, (Select count(1) from likes where post_id = posts.id AND value = 1) as likes, (Select count(1) from likes where post_id = posts.id AND value = -1) as dislikes, (select count(1) from posts as p where p.post_id = posts.id) as replies FROM posts JOIN users ON users.ID = posts.user_id WHERE post_id = ? ORDER BY ID DESC LIMIT 10")
		row, err = stmt.Query(pid)
	}

	for row.Next() {
		u := User{}
		p := Post{}
		var cr string

		row.Scan(&p.Id, &p.Content, &cr, &u.Name, &u.Username, &u.ProfileUrl, &p.Likes, &p.Dislikes, &p.Replies)
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

	stmt, err := conn.Prepare("SELECT posts.id, posts.content, posts.created_at, users.name, users.username, users.profile_url, (Select count(1) from likes where post_id = posts.id AND value = 1) as likes, (Select count(1) from likes where post_id = posts.id AND value = -1) as dislikes, (select count(1) from posts where post_id = posts.id) as replies FROM posts JOIN users ON users.ID = posts.user_id WHERE posts.id = ? ORDER BY ID DESC")
	row, _ := stmt.Query(pid)

	u := User{}
	p := Post{}
	for row.Next() {
		var cr string

		row.Scan(&p.Id, &p.Content, &cr, &u.Name, &u.Username, &u.ProfileUrl, &p.Likes, &p.Dislikes, &p.Replies)
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", cr)

		p.User = u
		p.Posts, _ = GetPosts(p.Id, 0)
	}

	return p
}

func AddPost(u User, content string, pid int) {
	conn, err := helpers.GetConn()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	c := content

	if pid == 0 {
		stmt, _ := conn.Prepare("INSERT INTO posts ( content, user_id, created_at ) VALUES ( ?, ?, NOW() )")
		stmt.Exec(c, u.Id)
	} else {
		stmt, _ := conn.Prepare("INSERT INTO posts ( content, user_id, created_at, post_id ) VALUES ( ?, ?, NOW(), ? )")
		stmt.Exec(c, u.Id, pid)
	}
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
