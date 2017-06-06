package models

import (
	"database/sql"
	"github.com/lib/pq"
	//	"strings"
	"errors"
	//"fmt"
	"strconv"
	"strings"
)

const ThreadTableCreationQuery = `CREATE TABLE IF NOT EXISTS thread
	(
		id SERIAL NOT NULL PRIMARY KEY,
		title VARCHAR(100),
		author VARCHAR(25) REFERENCES users(nickname) ON DELETE CASCADE ,
		forum VARCHAR(50) REFERENCES forum(slug) ON DELETE CASCADE,
		message TEXT NOT NULL ,
		votes INT NOT NULL DEFAULT 0,
		slug VARCHAR(25) NOT NULL,
		created TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);
	CREATE UNIQUE INDEX IF NOT EXISTS thread_slug_ci_index ON thread (lower(slug)) WHERE slug != '';
	CREATE INDEX IF NOT EXISTS thread_author_ci_index ON thread (lower(author));
	CREATE INDEX IF NOT EXISTS thread_forum_ci_index ON thread (lower(forum));
	CREATE UNIQUE INDEX IF NOT EXISTS thread_id_index ON thread (id);`

type Thread struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Votes   int    `json:"votes, omitempty"`
	Slug    string `json:"slug"`
	Created string `json:"created"`
}

func (t *Thread) ThreadCreateSQL(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.QueryRow(
		"INSERT INTO thread(title, author, forum, message, slug, created ) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, created",
		t.Title, t.Author, t.Forum, t.Message, t.Slug, t.Created).Scan(&t.ID, &t.Created)

	if err != nil {
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
	_, err = tx.Exec("UPDATE forum SET threads=threads+1 WHERE lower(slug)=$1", strings.ToLower(t.Forum))
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()

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
	return db.QueryRow("SELECT id, title, author, forum, message, slug, created FROM thread WHERE lower(slug)=$1",
		strings.ToLower(t.Slug)).Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Slug, &t.Created)
}

func (t *Thread) ThreadVoteSQL(db *sql.DB, voice int) error {
	return db.QueryRow("UPDATE thread SET votes = thread.votes + ($1)  WHERE id = $2 RETURNING votes", voice, t.ID).Scan(&t.Votes)
}

func (t *Thread) ThreadGetPostsFlatSQL(db *sql.DB, limit, desc string, offset int) ([]Post, error) {
	queryRow := `SELECT id, parent, author, message, isEdited, forum, thread, created FROM post WHERE thread=$1`

	var params []interface{}
	params = append(params, t.ID)
	paramOffset := 2

	if desc == "true" {
		queryRow += ` ORDER BY created DESC, id DESC`
	} else {
		queryRow += ` ORDER BY created ASC, id ASC`
	}
	if limit != "" {
		queryRow += ` LIMIT $` + strconv.Itoa(paramOffset)
		params = append(params, limit)
		paramOffset += 1
	}
	if offset != 0 {
		queryRow += ` OFFSET $` + strconv.Itoa(paramOffset)
		params = append(params, offset)
		paramOffset += 1
	}

	rows, err := db.Query(queryRow, params...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.Thread, &p.Created); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (t *Thread) ThreadGetPostsTreeSQL(db *sql.DB, limit, desc string, offset int) ([]Post, error) {
	queryRow := `SELECT id, parent, author, message, isEdited, forum, thread, created FROM post WHERE thread = $1`
//`	queryRow := `WITH RECURSIVE tree(id, parent, author, message, isEdited, forum, thread, created, path) AS(
//    SELECT id, parent, author, message, isEdited, forum, thread, created, ARRAY[id] FROM post WHERE thread = $1 AND parent=0
//    UNION
//      SELECT post.id, post.parent, post.author, post.message, post.isEdited, post.forum, post.thread, post.created, path||post.id FROM post
//         JOIN tree ON post.parent = tree.id
//      WHERE post.thread = $1
//) SELECT id, parent, author, message, isEdited, forum, thread, created FROM tree
//`

	var params []interface{}
	params = append(params, t.ID)
	paramOffset := 2

	if desc == "true" {
		queryRow += ` ORDER BY parentpath DESC`
	} else {
		queryRow += ` ORDER BY parentpath, created ASC`
	}
	if limit != "" {
		queryRow += ` LIMIT $` + strconv.Itoa(paramOffset)
		params = append(params, limit)
		paramOffset += 1
	}
	if offset != 0 {
		queryRow += ` OFFSET $` + strconv.Itoa(paramOffset)
		params = append(params, offset)
		paramOffset += 1
	}

	rows, err := db.Query(queryRow, params...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.Thread, &p.Created); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (t *Thread) ThreadGetPostsParentTreeSQL(db *sql.DB, limit, desc string, offset int) ([]Post, error) {
	queryRow := `SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
WHERE thread = $1 AND parentPath[1] in (
	SELECT id FROM post
	WHERE thread = $1 AND array_length(parentPath, 1) = 1`
//`WITH RECURSIVE tree(id, parent, author, message, isEdited, forum, thread, created, path) AS(
//    (SELECT id, parent, author, message, isEdited, forum, thread, created, ARRAY[id] FROM post WHERE thread = $1 AND parent=0`

	endQueryRow := ""
	var params []interface{}
	params = append(params, t.ID)
	paramOffset := 2

	if desc == "true" {
		queryRow += ` ORDER BY id DESC`
		endQueryRow += ` ORDER BY parentPath DESC;`
	} else {
		queryRow += ` ORDER BY id ASC`
		endQueryRow += ` ORDER BY parentPath ASC;`
	}
	if limit != "" {
		queryRow += ` LIMIT $` + strconv.Itoa(paramOffset)
		params = append(params, limit)
		paramOffset += 1
	}
	if offset != 0 {
		queryRow += ` OFFSET $` + strconv.Itoa(paramOffset)
		params = append(params, offset)
		paramOffset += 1
	}
	queryRow += `)`

	queryRow += endQueryRow

	//UNION
	//SELECT post.id, post.parent, post.author, post.message, post.isEdited, post.forum, post.thread, post.created, path||post.id FROM post
	//JOIN tree ON post.parent = tree.id
	//WHERE post.thread = $1
	//) SELECT id, parent, author, message, isEdited, forum, thread, created FROM tree` + endQueryRow

	rows, err := db.Query(queryRow, params...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.Thread, &p.Created); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (t *Thread) ThreadUpdateSQL(db *sql.DB) error {
	queryRow := "UPDATE thread SET "

	var params []interface{}
	paramOffset := 1

	if t.Message != "" {
		queryRow += ` message=$` + strconv.Itoa(paramOffset)
		params = append(params, t.Message)
		paramOffset += 1
	}
	if t.Title != "" {
		if paramOffset == 2 {
			queryRow += `,`
		}
		queryRow += ` title=$` + strconv.Itoa(paramOffset)
		params = append(params, t.Title)
		paramOffset += 1
	}
	if t.ID != 0 {
		queryRow += ` WHERE id=$` + strconv.Itoa(paramOffset)
		params = append(params, t.ID)
	} else if t.Slug != "" {
		queryRow += ` WHERE lower(slug)=$` + strconv.Itoa(paramOffset)
		params = append(params, strings.ToLower(t.Slug))
	} else {
		return errors.New("sasd")
	}

	queryRow += " RETURNING id, title, author, forum, message, slug, votes, created"
	err := db.QueryRow(queryRow, params...).Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Slug, &t.Votes, &t.Created)

	if err != nil {
		if err != sql.ErrNoRows {
			switch err.(*pq.Error).Code {
			case pq.ErrorCode("23505"):
				return UniqueError
			default:
				return err
			}
		}
		return err
	}

	return nil
}

func ThreadCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM thread").Scan(&count)
	return count, err
}
