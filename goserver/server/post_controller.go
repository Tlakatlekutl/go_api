package server

import (
	md "./models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SinglePostResponse struct {
	Sp     md.Post    `json:"post"`
	Ath    *md.User   `json:"author,omitempty,-"`
	Forum  *md.Forum  `json:"forum,omitempty"`
	Thread *md.Thread `json:"thread,omitempty"`
}

func PostGetOne(w http.ResponseWriter, r *http.Request) {
	var sid string = mux.Vars(r)["id"]
	id, _ := strconv.Atoi(sid)
	p := md.Post{Id: id}
	if err := p.PostGetOneSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	resp := SinglePostResponse{Sp: p}

	queryParams := r.URL.Query()
	if related, ok := queryParams["related"]; ok {
		params := strings.Split(related[0], ",")
		for _, task := range params {
			var err error
			switch task {
			case "user":
				u := md.User{Nickname: p.Author}
				err = u.GetOneUserSQL(DB.DB)
				resp.Ath = &u

			case "forum":
				f := md.Forum{Slug: p.Forum}
				err = f.GetForumByUniqueSlug(DB.DB)
				resp.Forum = &f
			case "thread":
				t := md.Thread{ID: p.Thread}
				err = t.ThreadSelectOneIdOrSlugSQL(DB.DB)
				resp.Thread = &t
			}
			if err != nil {
				CheckDbErr(err, w)
				return
			}

		}
	}

	RespondJSON(w, http.StatusOK, resp)
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {
	var sid string = mux.Vars(r)["id"]
	id, _ := strconv.Atoi(sid)
	p := md.Post{Id: id}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid parse post update json")
		return
	}
	defer r.Body.Close()

	if p.Message == "" {
		if err := p.PostGetOneSQL(DB.DB); err != nil {
			CheckDbErr(err, w)
			return
		}
		RespondJSON(w, http.StatusOK, p)
		return
	}

	if err := p.PostUpdateSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	RespondJSON(w, http.StatusOK, p)
}

func PostsCreate(w http.ResponseWriter, r *http.Request) {
	var ppk string = mux.Vars(r)["slug_or_id"]
	pa := []md.Post{}
	t := md.Thread{}
	created := time.Now()
	timeStamp := created.Format(time.RFC3339)

	if id, err := IsId(ppk); err == nil {
		t.ID = id
	} else {
		t.Slug = ppk
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pa); err != nil {

		RespondError(w, http.StatusBadRequest, "Invalid request creation post")
		return
	}
	defer r.Body.Close()

	if err := t.ThreadSelectOneIdOrSlugSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	if err := md.PostCreateListSQL(DB.DB, pa, t.Forum, timeStamp, t.ID); err != nil {
		CheckDbErr(err, w)
		return
	}


	RespondJSON(w, http.StatusCreated, pa)
}
