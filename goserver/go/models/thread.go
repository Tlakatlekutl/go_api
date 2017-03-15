package models

import (
	"database/sql"
	"github.com/lib/pq"
//	"strings"
	"errors"
	//"fmt"
	"strings"
	"strconv"
)

const ThreadTableCreationQuery =
	`CREATE TABLE IF NOT EXISTS thread
	(
		id SERIAL NOT NULL PRIMARY KEY,
		title VARCHAR(100),
		author VARCHAR(25) REFERENCES users(nickname),
		forum VARCHAR(50) REFERENCES forum(slug),
		message TEXT NOT NULL ,
		votes INT NOT NULL DEFAULT 0,
		slug VARCHAR(25) NOT NULL,
		created TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);
	CREATE UNIQUE INDEX IF NOT EXISTS thread_slug_ci_index ON thread (lower(slug)) WHERE slug != '';
-- 	--CREATE UNIQUE INDEX IF NOT EXISTS user_nickname_ci_index ON users ((lower(nickname)));`


type Thread struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Forum string `json:"forum"`
	Message string `json:"message"`
	Votes int `json:"votes, omitempty"`
	Slug string `json:"slug"`
	Created string `json:"created"`
}

func (t *Thread) ThreadCreateSQL(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO thread(title, author, forum, message, slug, created ) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		t.Title, t.Author, t.Forum, t.Message, t.Slug, t.Created).Scan(&t.ID)
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

func (t *Thread) ThreadSelectOneIdOrSlugSQL(db *sql.DB) error {
	if t.ID != 0 {
		return db.QueryRow("SELECT  title, author, forum, message, slug, votes, created FROM thread WHERE id=$1",
			t.ID).Scan(&t.Title, &t.Author, &t.Forum, &t.Message, &t.Slug, &t.Votes, &t.Created)
	} else if t.Slug != "" {
		return db.QueryRow("SELECT id, title, author, forum, message, slug, votes, created FROM thread WHERE lower(slug)=$1",
			strings.ToLower(t.Slug)).Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Slug, &t.Votes, &t.Created)
	} else {
		return errors.New("something was wrong Thread")
	}
}

func (t *Thread) ThreadGetOneSQL(db *sql.DB) error {
	return db.QueryRow("SELECT id, title, author, forum, message, slug,created FROM thread WHERE lower(slug)=$1",
		strings.ToLower(t.Slug)).Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Slug, &t.Created)
}
func (t *Thread) ThreadVoteSQL(db *sql.DB, voice int) error {
	return db.QueryRow("UPDATE thread SET votes = thread.votes + ($1)  WHERE id = $2 RETURNING votes", voice, t.ID).Scan(&t.Votes)
}

func (t *Thread) ThreadGetListsPostsSQL(db *sql.DB, limit, since, desc string) ([]Post, error) {
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

