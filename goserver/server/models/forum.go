package models

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"strings"
	//"fmt"
	"strconv"
)

const ForumTableCreationQuery =
	`CREATE TABLE IF NOT EXISTS forum
	(
		title VARCHAR(200) NOT NULL,
		userPK VARCHAR(25) REFERENCES users(nickname) ON DELETE CASCADE,
		slug VARCHAR(50) PRIMARY KEY,
		posts INT DEFAULT 0,
		threads INT DEFAULT 0
	);
	CREATE UNIQUE INDEX IF NOT EXISTS forum_slug_ci_index ON forum ((lower(slug)));
	CREATE UNIQUE INDEX IF NOT EXISTS forum_pk_index ON forum (lower(userPK));`

var FKConstraintError = errors.New("violates foreign key constraint")

type Forum struct {
	id int
	Title string `json:"title"`
	User string `json:"user"`
	Slug string `json:"slug"`
	Posts int `json:"posts"`
	Threads int `json:"threads"`
}

func (f *Forum) ForumCreateSQL(db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO forum(title, slug, userPK) VALUES($1, $2, $3)",
		f.Title, f.Slug, f.User)
	if err!=nil {
		switch err.(*pq.Error).Code {
		case pq.ErrorCode("23505"):
			return UniqueError
		case pq.ErrorCode("23503"):
			return FKConstraintError
		default:
			return err
		}
	}
	return err
}

func (f *Forum) GetForumByUniqueNickname(db * sql.DB) error {
	return db.QueryRow("SELECT title, slug, posts, threads FROM forum WHERE lower(userPK)=$1",
		strings.ToLower(f.User)).Scan(&f.Title, &f.Slug, &f.Posts, &f.Threads)
}

func (f *Forum) GetForumByUniqueSlug(db * sql.DB) error {
	return db.QueryRow("SELECT title, userPk, slug, posts, threads FROM forum WHERE lower(slug)=$1",
		strings.ToLower(f.Slug)).Scan(&f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads)
}

func (f *Forum) ForumGetListThreadsSQL(db *sql.DB, limit, since, desc string) ([]Thread, error) {
	queryRow := `SELECT id, title, author, forum, message, votes, slug, created FROM thread WHERE lower(forum)=$1`

	var params []interface{}
	params = append(params, strings.ToLower(f.Slug))
	paramOffset := 2
	if since != "" && desc == "true"{
		queryRow += ` AND created <= $`+strconv.Itoa(paramOffset)
		params = append(params, since)
		paramOffset+=1
	} else if since !="" {
		queryRow += ` AND created >= $`+strconv.Itoa(paramOffset)
		params = append(params, since)
		paramOffset+=1
	}
	if desc == "true" {
		queryRow += ` ORDER BY created DESC`
	} else {
		queryRow += ` ORDER BY created ASC`
	}
	if limit != "" {
		queryRow += ` LIMIT $`+strconv.Itoa(paramOffset)
		params = append(params, limit)
		paramOffset+=1
	}


	rows, err := db.Query(queryRow, params...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := []Thread{}
	for rows.Next() {
		var t Thread
		if err := rows.Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message,&t.Votes, &t.Slug, &t.Created); err!=nil {
			return nil, err
		}
		threads = append(threads, t)
	}

	return threads, nil
}

func (f *Forum) ForumGetListUsersSQL(db *sql.DB, limit, since, desc string) ([]User, error) {
	queryRow := `SELECT DISTINCT lower(nickname) COLLATE "ucs_basic", nickname, fullname, about, email FROM users u
	LEFT JOIN thread t ON t.author=u.nickname
	LEFT JOIN post p ON p.author=u.nickname
	LEFT JOIN forum f ON f.userpk=u.nickname
	WHERE (lower(t.forum)=$1 OR lower(p.forum) = $1 OR lower(f.userpk)=$1)`

	var params []interface{}
	params = append(params, strings.ToLower(f.Slug))

	paramOffset := 2
	if since != "" && desc == "true"{
		queryRow += ` AND lower(nickname) COLLATE "ucs_basic" < lower($`+strconv.Itoa(paramOffset)+`) COLLATE "ucs_basic"`
		params = append(params, since)
		paramOffset+=1
	} else if since !="" {
		queryRow += ` AND lower(nickname) COLLATE "ucs_basic" > lower($`+strconv.Itoa(paramOffset)+`) COLLATE "ucs_basic"`
		params = append(params, since)
		paramOffset+=1
	}

	if desc == "true" {
		queryRow += ` ORDER BY lower(nickname) COLLATE "ucs_basic" DESC`
	} else {
		queryRow += ` ORDER BY lower(nickname) COLLATE "ucs_basic" ASC`
	}
	if limit != "" {
		queryRow += ` LIMIT $`+strconv.Itoa(paramOffset)
		params = append(params, limit)
		paramOffset+=1
	}

	rows, err := db.Query(queryRow, params...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		var temp string
		if err := rows.Scan(&temp, &u.Nickname, &u.Fullname, &u.About, &u.Email); err!=nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func ForumCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM forum").Scan(&count)
	return count, err
}
