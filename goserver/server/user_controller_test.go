package server

import (
	md "./models"
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"encoding/json"
	"bytes"
	"net/http/httptest"
	"github.com/lib/pq"
	"strings"
)

func TestApp_UserCreateSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error '%s' during db mock", err)
	}
	defer db.Close()

	a := App{DB: db}

	u := md.User{
		Nickname:"SomeNickname1",
		Email:"test@golang.com",
		About: "Have fun",
		Fullname: "Gofer",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(u.Nickname, u.Fullname, u.Email, u.About).
		WillReturnResult(sqlmock.NewResult(1,1))


	body, _ := json.Marshal(u)

	req, err := http.NewRequest("POST", "/api/user/SomeNickname1/create", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("NewRequest error '%s' ", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(a.UserCreate)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Wrong status: expected 200, got %v", rr.Code)
	}
}

func TestApp_UserCreateUniqueError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error '%s' during db mock", err)
	}
	defer db.Close()

	a := App{DB: db}

	u := md.User{
		Nickname:"SomeNickname1",
		Email:"test@golang.com",
		About: "Have fun",
		Fullname: "Gofer",
	}
	pqerr := &pq.Error{Code:"23505"}
	mock.ExpectExec("INSERT INTO users").
		WithArgs(u.Nickname, u.Fullname, u.Email, u.About).
		WillReturnError(pqerr)


	u2 := md.User{
		Nickname:"2SomeNickname1",
		Email:"2test@golang.com",
		About: "2Have fun",
		Fullname: "2Gofer",
	}
	rows := sqlmock.NewRows([]string{"nickname", "fullname", "email", "about"}).
		AddRow(u2.Nickname, u2.Fullname, u2.Email, u2.About)
	mock.ExpectQuery("SELECT nickname, fullname, email, about FROM users").
		WithArgs(strings.ToLower(u.Nickname), strings.ToLower(u.Email)).
		WillReturnRows(rows)


	body, _ := json.Marshal(u)

	req, err := http.NewRequest("POST", "/api/user/SomeNickname1/create", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("NewRequest error '%s' ", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(a.UserCreate)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("Wrong status: expected 409, got %v", rr.Code)
	}
}
