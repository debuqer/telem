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
}

type Posts []Post

func GetFeed() ([]Post, error) {
	posts := Posts{}

	conn, err := helpers.GetConn()
	if err != nil {
		return make(Posts, 0), err
	}

	stmt, err := conn.Prepare("SELECT posts.id, posts.content, posts.created_at, users.name, users.username, users.profile_url FROM posts JOIN users ON users.ID = posts.user_id ORDER BY ID DESC LIMIT 10")
	row, _ := stmt.Query()

	for row.Next() {
		u := User{}
		p := Post{}
		var cr string

		row.Scan(&p.Id, &p.Content, &cr, &u.Name, &u.Username, &u.ProfileUrl)
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", cr)

		p.User = u

		posts = append(posts, p)
	}

	return posts, nil
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
