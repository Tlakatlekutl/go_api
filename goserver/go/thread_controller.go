package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	md "./models"
	//"fmt"
)


func ThreadCreate(w http.ResponseWriter, r *http.Request) {
	var fsg string = mux.Vars(r)["slug"]
	f := md.Forum{Slug: fsg}

	if err := f.GetForumByUniqueSlug(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	t := md.Thread{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request creation thread")
		return
	}
	defer r.Body.Close()

	t.Forum = f.Slug

	if t.Created == "" {
		t.Created="2017-08-22T01:30:51.934+03:00"
	}

	if err := t.ThreadCreateSQL(DB.DB); err != nil {
		switch err {
		case md.UniqueError:
			if  err := t.ThreadGetOneSQL(DB.DB); err == nil {
				RespondJSON(w, http.StatusConflict, t)
				return
			}
		case md.FKConstraintError:
			RespondError(w, http.StatusNotFound, err.Error())

		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, http.StatusCreated, t)
}

func ThreadGetOne(w http.ResponseWriter, r *http.Request) {
	var ppk string = mux.Vars(r)["slug_or_id"]
	t:=md.Thread{}

	if id, err := IsId(ppk); err==nil {
		t.ID = id
	} else {
		t.Slug = ppk
	}
	if err := t.ThreadSelectOneIdOrSlugSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}
	RespondJSON(w, http.StatusOK, t)

}

func ThreadGetPosts(w http.ResponseWriter, r *http.Request) {
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

func ThreadUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadVote(w http.ResponseWriter, r *http.Request) {
	var ppk string = mux.Vars(r)["slug_or_id"]
	t:=md.Thread{}

	if id, err := IsId(ppk); err==nil {
		t.ID = id
	} else {
		t.Slug = ppk
	}
	if err := t.ThreadSelectOneIdOrSlugSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	v := md.Vote{Thread: t.ID}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid parse vote json")
		return
	}
	defer r.Body.Close()

	u := md.User{Nickname: v.Nickname}

	if err := u.GetOneUserSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	if err := v.VoteSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}
	if sum, err := v.VoteCountSQL(DB.DB); err == nil {
		t.Votes = sum
	} else {
		CheckDbErr(err, w)
		return
	}

	RespondJSON(w, http.StatusOK, t)

}
