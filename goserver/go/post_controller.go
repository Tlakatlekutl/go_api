package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	md "./models"
	"strconv"
	"strings"
)

type SinglePostResponse struct {
	Sp md.Post `json:"post"`
	Ath *md.User `json:"author,omitempty,-"`
	Forum *md.Forum `json:"forum,omitempty"`
	Thread *md.Thread `json:"thread,omitempty"`
}

func PostGetOne(w http.ResponseWriter, r *http.Request) {
	var sid string = mux.Vars(r)["id"]
	id, _ := strconv.Atoi(sid)
	p := md.Post{Id:id}
	if err := p.PostGetOneSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	resp := SinglePostResponse{Sp:p}

	queryParams := r.URL.Query()
	if related, ok := queryParams["related"]; ok {
		params := strings.Split(related[0], ",")
		for _, task :=range params {
			var err error
			switch task {
			case "user":
				u:=md.User{Nickname: p.Author}
				err = u.GetOneUserSQL(DB.DB)
				resp.Ath=&u

			case "forum":
				f:=md.Forum{Slug: p.Forum}
				err = f.GetForumByUniqueSlug(DB.DB)
				resp.Forum=&f
			case "thread":
				t:=md.Thread{ID: p.Thread}
				err = t.ThreadSelectOneIdOrSlugSQL(DB.DB)
				resp.Thread=&t
			}
			if err!= nil {
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
	p := md.Post{Id:id}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid parse post update json")
		return
	}
	defer r.Body.Close()

	if p.Message=="" {
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
	t:=md.Thread{}

	if id, err := IsId(ppk); err==nil {
		t.ID = id
	} else {
		t.Slug = ppk
	}
	//fmt.Println(t.ID, "slog1:", t.Slug)

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
	for i:=0; i < len(pa); i+=1 {
		pa[i].Thread = t.ID
		pa[i].Forum = t.Forum

		if err := pa[i].PostCreateOneSQL(DB.DB); err != nil {
			CheckDbErr(err, w)
			return
		}
	}

	RespondJSON(w, http.StatusCreated, pa)
}

