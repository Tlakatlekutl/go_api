package models

import (
	"database/sql"
	"github.com/lib/pq"
)

const ThreadTableCreationQuery =
	`CREATE TABLE IF NOT EXISTS thread
	(
		id SERIAL NOT NULL PRIMARY KEY,
		title VARCHAR(80),
		author VARCHAR(25) REFERENCES users(nickname),
		forum VARCHAR(50) REFERENCES forum(slug),
		message TEXT NOT NULL ,
		votes INT NOT NULL DEFAULT 0,
		slug VARCHAR(25) NOT NULL,
		created TIMESTAMP NOT NULL DEFAULT current_timestamp
	);
	CREATE UNIQUE INDEX IF NOT EXISTS forum_slug_ci_index ON forum ((lower(slug)));
-- 	CREATE UNIQUE INDEX IF NOT EXISTS user_nickname_ci_index ON users ((lower(nickname)));`


type Thread struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Forum string `json:"forum"`
	Message string `json:"message"`
	Votes int `json:"votes"`
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