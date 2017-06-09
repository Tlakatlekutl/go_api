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

//
//INSERT INTO distributors (did, dname)
//VALUES (5, 'Gizmo Transglobal'), (6, 'Associated Computing, Inc')
//ON CONFLICT (did) DO UPDATE SET dname = EXCLUDED.dname;
