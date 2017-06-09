package models

import (
	"database/sql"
)

type ForumUser struct {
	Forum string
	UserPK string
}

const ForumUsersTableCreationQuery = `CREATE TABLE IF NOT EXISTS forum_user
	(
		forum VARCHAR(50) REFERENCES forum(slug) ON DELETE CASCADE,
		userPK VARCHAR(25) REFERENCES users(nickname) ON DELETE CASCADE
	);
	ALTER TABLE forum_user DROP CONSTRAINT IF EXISTS unique_pair_constr_fu;
	ALTER TABLE forum_user ADD CONSTRAINT  unique_pair_constr_fu UNIQUE (forum, userPK);
	CREATE INDEX IF NOT EXISTS fu_forum_index on forum_user(lower(forum));
	CREATE INDEX IF NOT EXISTS fu_user_index on forum_user(lower(userPK));
	CREATE INDEX IF NOT EXISTS fu_user_collate_index on forum_user(lower(userPK) COLLATE "ucs_basic");
	`

func (fu *ForumUser) ForumUserInsertSQL(db *sql.DB) error {
	if _, err := db.Exec(
		`INSERT INTO forum_user(forum, userPK) VALUES($1, $2)
ON CONFLICT ON CONSTRAINT unique_pair_constr_fu DO NOTHING`,
		fu.Forum, fu.UserPK);
		err != nil {
			return err
	}
	return nil
}