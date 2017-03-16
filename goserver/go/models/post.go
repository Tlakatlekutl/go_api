package models

import (
	"database/sql"
	//"errors"
	"github.com/lib/pq"
	//"strings"
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
	err := db.QueryRow(
		//"INSERT INTO post(parent, author, message, isedited, forum, thread, created) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		//p.Parent, p.Author, p.Message, p.IsEdited, p.Forum, p.Thread, p.Created).Scan(&p.Id)
		"INSERT INTO post(parent, author, message, isedited, forum, thread) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, created",
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
	}
	return err
}