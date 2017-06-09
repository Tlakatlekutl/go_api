package models

import (
	"database/sql"
	//"strings"
	"log"
	"github.com/lib/pq"
	"strings"
	"fmt"
)

const PostTableCreationQuery = `CREATE TABLE IF NOT EXISTS post
	(
		id SERIAL NOT NULL PRIMARY KEY,
		parent INT NOT NULL DEFAULT 0,
		author VARCHAR(25) NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
		message TEXT NOT NULL ,
		isEdited BOOLEAN DEFAULT FALSE ,
		forum VARCHAR(50) NOT NULL REFERENCES forum(slug) ON DELETE CASCADE,
		thread INT NOT NULL REFERENCES thread(id) ON DELETE CASCADE,
		created TIMESTAMPTZ,
		parentPath INT[]
	);
 	CREATE INDEX IF NOT EXISTS post_author_ci_index ON post((lower(author)));
 	CREATE INDEX IF NOT EXISTS post_forum_ci_index ON post ((lower(forum)));
 	CREATE INDEX IF NOT EXISTS post_thread_ci_index ON post (thread);
 	CREATE UNIQUE INDEX IF NOT EXISTS post_id_parent_index ON post (id, parent);
 	CREATE INDEX IF NOT EXISTS post_thread_id ON post (thread, id);
 	CREATE UNIQUE INDEX IF NOT EXISTS post_id_index ON post(id);
 	CREATE INDEX IF NOT EXISTS parent_path_second_elem_index on post ((parentPath[1]));
 	`

type Post struct {
	Id       int    `json:"id"`
	Parent   int    `json:"parent"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited"`
	Forum    string `json:"forum"`
	Thread   int    `json:"thread"`
	Created  string `json:"created"`
}


func PostCreateListSQL(db *sql.DB, postList []Post, forum, created string, thread int) error {

	tx, err := db.Begin()
	//defer tx.Rollback()
	if err != nil {
		log.Fatal(err);
	}

	stmt, err := tx.Prepare(pq.CopyIn("post", "id","parent", "author", "message", "isedited", "forum", "thread", "created","parentpath"))

	uniqueUsers := []ForumUser{}
	for i:=0; i < len(postList); i+=1 {
		err = db.QueryRow("SELECT nextval(pg_get_serial_sequence('post', 'id'))").Scan(&postList[i].Id)
		if err != nil {
			tx.Rollback()
			return err
		}
		//parentPath := []sql.NullInt64{}
		parentPath := []int64{}
		if (postList[i].Parent != 0) {
			//var path []sql.NullInt64{}
			var path []int64
			err = db.QueryRow("SELECT parentPath FROM post WHERE id=$1 AND thread=$2", postList[i].Parent, thread).Scan(pq.Array(&path))
			if err != nil {
				//tx.Rollback()
				//return err
				return UniqueError
			}
			parentPath = append(parentPath, path...)
		}
		parentPath = append(parentPath, int64(postList[i].Id))
		//parentPath = append([]int64{int64(postList[i].Parent)}, parentPath...)
		postList[i].Forum = forum
		postList[i].Created = created
		postList[i].Thread = thread

		_, err  = stmt.Exec(postList[i].Id, postList[i].Parent, postList[i].Author,postList[i].Message, postList[i].IsEdited, forum, thread, created, pq.Array(parentPath));
		if err != nil {
			//tx.Rollback()
			return err;
		}
		uniqueUsers = UniqArray(uniqueUsers, postList[i].Author, forum)
	}
	_ , err = stmt.Exec()

	if err != nil {
		tx.Rollback()
		return parseError(err)
	}

	err = stmt.Close()
	if err != nil {
		tx.Rollback()
		return parseError(err)
	}


	_, err = tx.Exec("UPDATE forum SET posts=posts+$1 WHERE lower(slug)=$2",len(postList), strings.ToLower(forum))
	if err != nil {
		tx.Rollback()
		return parseError(err)
	}
	for _, fu := range uniqueUsers {
		_, err = tx.Exec(
			`INSERT INTO forum_user(forum, userPK) VALUES($1, $2)
			ON CONFLICT ON CONSTRAINT unique_pair_constr_fu DO NOTHING`,
			forum, fu.UserPK);
		if (parseError(err) == UniqueError) {
			fmt.Println(err.Error())
			err = nil;
		} else if err != nil{
			fmt.Println(err.Error())
			tx.Rollback()
			return parseError(err)
		}
	}

	tx.Commit();

	return parseError(err);

}

func (p *Post) PostGetOneSQL(db *sql.DB) error {
	return db.QueryRow("SELECT parent, author, message, isedited, forum, thread, created  FROM post WHERE id=$1", p.Id).Scan(
		&p.Parent, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.Thread, &p.Created)
}

func (p *Post) PostUpdateSQL(db *sql.DB) error {
	return db.QueryRow("UPDATE post SET message=$2, isedited=(message != $2) WHERE id=$1 RETURNING parent, author, isedited, forum, thread, created", p.Id, p.Message).Scan(
		&p.Parent, &p.Author, &p.IsEdited, &p.Forum, &p.Thread, &p.Created)
}

func PostCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM post").Scan(&count)
	return count, err
}
