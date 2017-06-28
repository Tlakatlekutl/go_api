package server

import (
	md "./models"
	"net/http"
)

type StatusResponse struct {
	Forum  int `json:"forum"`
	Post   int `json:"post"`
	Thread int `json:"thread"`
	User   int `json:"user"`
}

func (a *App) Status(w http.ResponseWriter, r *http.Request) {
	var res StatusResponse
	var err error
	if res.Forum, err = md.ForumCount(a.DB); err != nil {
		CheckDbErr(err, w)
		return
	}
	if res.Post, err = md.PostCount(a.DB); err != nil {
		CheckDbErr(err, w)
		return
	}
	if res.Thread, err = md.ThreadCount(a.DB); err != nil {
		CheckDbErr(err, w)
		return
	}
	if res.User, err = md.UserCount(a.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	RespondJSON(w, http.StatusOK, res)

}

func (a *App) Clear(w http.ResponseWriter, r *http.Request) {
	if _, err := a.DB.Exec("DELETE FROM vote"); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
	}
	if _, err := a.DB.Exec("DELETE FROM post"); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
	}
	if _, err := a.DB.Exec("DELETE FROM thread"); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
	}
	if _, err := a.DB.Exec("DELETE FROM forum"); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
	}
	if _, err := a.DB.Exec("DELETE FROM users"); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
	}
	RespondJSON(w, http.StatusOK, nil)

}
