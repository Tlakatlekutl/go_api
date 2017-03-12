package server

import (
	"net/http"
	//"github.com/gorilla/mux"
	"encoding/json"
	md "./models"
	"github.com/gorilla/mux"
	//"fmt"
)

func Clear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

}

func ForumCreate(w http.ResponseWriter, r *http.Request) {
	f := md.Forum{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&f); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request creation forum")
		return
	}
	defer r.Body.Close()

	u := md.User{Nickname: f.User}

	if err := u.GetOneUserSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	f.User = u.Nickname

	if err := f.ForumCreateSQL(DB.DB); err != nil {
		switch err {
		case md.UniqueError:
			if err := f.GetForumByUniqueNickname(DB.DB); err == nil {
				RespondJSON(w, http.StatusConflict, f)
				return
			}
		case md.FKConstraintError:
			RespondError(w, http.StatusNotFound, err.Error())

		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, http.StatusCreated, f)
}

func ForumGetOne(w http.ResponseWriter, r *http.Request) {
	var sg string = mux.Vars(r)["slug"]
	f := md.Forum{Slug: sg}

	if err := f.GetForumByUniqueSlug(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	RespondJSON(w, http.StatusOK, f)
}

func ForumGetThreads(w http.ResponseWriter, r *http.Request) {
	var sg string = mux.Vars(r)["slug"]

	queryParams := r.URL.Query()
	var limit string
	if val, ok := queryParams["limit"]; ok {
		limit = val[0]
	}
	var since string
	if val, ok := queryParams["since"]; ok {
		since = val[0]
	}
	var desc string
	if val, ok := queryParams["desc"]; ok {
		desc = val[0]
	}

	//fmt.Println(limit, since, desc)

	f := md.Forum{Slug: sg}
	if err := f.GetForumByUniqueSlug(DB.DB); err!=nil {
		CheckDbErr(err, w)
		return
	}
	if threads, err := f.ForumGetListThreadsSQL(DB.DB, limit, since, desc); err != nil {
		CheckDbErr(err, w)
		return
	} else {
		RespondJSON(w, http.StatusOK, threads)
	}
}
