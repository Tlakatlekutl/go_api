package models

import (
	"database/sql"
	"github.com/lib/pq"
)

const VoteTableCreationQuery = `CREATE TABLE IF NOT EXISTS vote
	(
		id SERIAL NOT NULL PRIMARY KEY,
		userPK VARCHAR(25) REFERENCES users(nickname),
		voice SMALLINT NOT NULL,
		thread INT REFERENCES thread(id)
		--changed BOOLEAN DEFAULT FALSE
	);
	ALTER TABLE vote DROP CONSTRAINT IF EXISTS unique_pair_constr;
	ALTER TABLE vote ADD CONSTRAINT unique_pair_constr UNIQUE (userPK, thread);
	--CREATE UNIQUE INDEX IF NOT EXISTS unique_pair ON vote (thread, userPK);	`

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
	Thread   int
}

func (v *Vote) VoteSQL(db *sql.DB) error {
	_, err := db.Exec(
		`INSERT INTO vote(userpk, voice, thread) VALUES($1, $2, $3)
		ON CONFLICT ON CONSTRAINT unique_pair_constr DO UPDATE SET voice = EXCLUDED.voice`, v.Nickname, v.Voice, v.Thread)
	if err != nil {
		switch err.(*pq.Error).Code {
		case pq.ErrorCode("23505"):
			return UniqueError
		default:
			return err
		}
	}

	return nil
}

func (v *Vote) VoteCountSQL(db *sql.DB) (int, error) {
	var sum int
	err := db.QueryRow("UPDATE thread SET votes = (SELECT SUM(voice) FROM vote WHERE thread=$1 GROUP BY thread) WHERE id = $1 RETURNING votes", v.Thread).Scan(&sum)
	return sum, err
}

func (t *Thread) ThreadVote(db *sql.DB, v *Vote) error {
	var vote int

	_, err := db.Exec(
		`INSERT INTO vote(userpk, voice, thread) VALUES($1, $2, $3)`,
		v.Nickname, v.Voice, v.Thread)
	if err == nil {
		err = db.QueryRow(`UPDATE thread SET votes = votes + $2 WHERE id=$1 RETURNING title, author, forum, message, slug, votes, created`, v.Thread, v.Voice).Scan(&t.Title, &t.Author, &t.Forum, &t.Message, &t.Slug, &t.Votes, &t.Created)
		if err != nil {
			//tx.Rollback()
			return parseError(err)
		}
	} else if parseError(err) == UniqueError {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		err = tx.QueryRow(`SELECT voice FROM vote WHERE userPK=$1 AND thread=$2`, v.Nickname, v.Thread).Scan(&vote)
		if (err != nil) {
			tx.Rollback()
			return parseError(err)
		}
		_, err = tx.Exec(
			`UPDATE vote SET voice = $3 WHERE userPK=$1 AND thread=$2`,
			v.Nickname, v.Thread, v.Voice)
		if (err != nil) {
			tx.Rollback()
			return parseError(err)
		}
		err = tx.QueryRow(`UPDATE thread SET votes = votes + $2 WHERE id=$1 RETURNING title, author, forum, message, slug, votes, created`, v.Thread, v.Voice-vote).Scan(&t.Title, &t.Author, &t.Forum, &t.Message, &t.Slug, &t.Votes, &t.Created)
		if (err != nil) {
			tx.Rollback()
			return parseError(err)
		}
		tx.Commit()
	} else {
		return parseError(err)
	}

	return nil

}
