package models

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"strings"
)

const UserTableCreationQuery = `CREATE TABLE IF NOT EXISTS users
	(
		nickname VARCHAR(25) PRIMARY KEY,
		fullname VARCHAR(60),
		email VARCHAR(50) NOT NULL UNIQUE,
		about TEXT
	);
	CREATE UNIQUE INDEX IF NOT EXISTS user_email_ci_index ON users ((lower(email)));
	CREATE UNIQUE INDEX IF NOT EXISTS user_nickname_ci_index ON users ((lower(nickname)));
	CREATE INDEX IF NOT EXISTS user_user_collate_index on users(lower(nickname) COLLATE "ucs_basic");`

var UniqueError = errors.New("unique")

type User struct {
	//id int
	Nickname string `json:"nickname"`
	About    string `json:"about"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

func (u *User) GetOneUserSQL(db *sql.DB) error {
	return db.QueryRow("SELECT nickname, fullname, email, about FROM users WHERE lower(nickname)=$1",
		strings.ToLower(u.Nickname)).Scan(&u.Nickname, &u.Fullname, &u.Email, &u.About)
}

func (u *User) CreateUserSQL(db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO users(nickname, fullname, email, about) VALUES($1, $2, $3, $4)",
		u.Nickname, u.Fullname, u.Email, u.About)
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

func (u *User) GetUniqueUsersSQL(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT nickname, fullname, email, about FROM users WHERE lower(nickname)=$1 OR lower(email) = $2",
		strings.ToLower(u.Nickname), strings.ToLower(u.Email))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Nickname, &u.Fullname, &u.Email, &u.About); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (u *User) UpdateUserSQL(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE users SET fullname=$1, email=$2, about=$3 WHERE lower(nickname)=$4",
		u.Fullname, u.Email, u.About, strings.ToLower(u.Nickname))

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

func (u *User) GetUsersListSQL(db *sql.DB, start, count int) ([]User, error) {
	return nil, errors.New("Not implemented")
}

func UserCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}
