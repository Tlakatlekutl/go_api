package server

import (
	"net/http"
	md"./models"
)


type StatusResponse struct {
	Forum int `json:"forum"`
	Post int `json:"post"`
	Thread int `json:"thread"`
	User int `json:"user"`
}

func Status(w http.ResponseWriter, r *http.Request) {
	var res StatusResponse
	var err error
	if res.Forum, err = md.ForumCount(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}
	if res.Post, err = md.PostCount(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}
	if res.Thread, err = md.ThreadCount(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}
	if res.User, err = md.UserCount(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	RespondJSON(w, http.StatusOK, res)

}



