package models

import (
	"database/sql"
	"strings"
//	"github.com/lib/pq"
	"github.com/lib/pq"
)

const PostTableCreationQuery =
	`CREATE TABLE IF NOT EXISTS post
	(
		id SERIAL NOT NULL PRIMARY KEY,
		parent INT NOT NULL DEFAULT 0,
		author VARCHAR(25) NOT NULL REFERENCES users(nickname),
		message TEXT NOT NULL ,
		isEdited BOOLEAN DEFAULT FALSE ,
		forum VARCHAR(50) NOT NULL REFERENCES forum(slug),
		thread INT NOT NULL REFERENCES thread(id),
		created TIMESTAMPTZ DEFAULT current_timestamp
	);
-- 	CREATE UNIQUE INDEX IF NOT EXISTS user_email_ci_index ON users ((lower(email)));`

type Post struct {
	Id int `json:"id"`
	Parent int `json:"parent"`
	Author string `json:"author"`
	Message string `json:"message"`
	IsEdited bool `json:"isEdited"`
	Forum string `json:"forum"`
	Thread int `json:"thread"`
	Created string `json:"created"`
}

func (p *Post) PostCreateOneSQL(db *sql.DB) error  {
	tx, err := db.Begin()
	if err!=nil {
		tx.Rollback()
		return err
	}
	var temp int
	if p.Parent!=0 {
		err = tx.QueryRow("SELECT id FROM post WHERE id=$1 AND thread=$2", p.Parent, p.Thread).Scan(&temp)
		if err!=nil {
			tx.Rollback()
			return UniqueError
		}
	}
	err = tx.QueryRow(
		`INSERT INTO post(parent, author, message, isedited, forum, thread) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, created;`,
	p.Parent, p.Author, p.Message, p.IsEdited, p.Forum, p.Thread).Scan(&p.Id, &p.Created)
	if err!=nil {
		switch err.(*pq.Error).Code {
		case pq.ErrorCode("23505"):
			return UniqueError
		case pq.ErrorCode("23503"):
			return FKConstraintError
		default:
			return err
		}
		tx.Rollback()
		return err
	}
	_, err =tx.Exec("UPDATE forum SET posts=posts+1 WHERE lower(slug)=$1", strings.ToLower(p.Forum))
	if err!=nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()

	return err
}

func (p *Post)PostGetOneSQL(db *sql.DB) error {
	return db.QueryRow("SELECT parent, author, message, isedited, forum, thread, created  FROM post WHERE id=$1", p.Id).Scan(
		&p.Parent, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.Thread, &p.Created)
}

func (p *Post)PostUpdateSQL(db *sql.DB) error {
	return db.QueryRow("UPDATE post SET message=$2, isedited=TRUE WHERE id=$1 RETURNING parent, author, isedited, forum, thread, created", p.Id, p.Message).Scan(
		&p.Parent, &p.Author, &p.IsEdited, &p.Forum, &p.Thread, &p.Created)
}

func PostCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM post").Scan(&count)
	return count, err
}