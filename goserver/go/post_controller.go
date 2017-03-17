package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	md "./models"
	//"fmt"
	"strconv"
)

type SinglePostResponse struct {
	Sp md.Post `json:"post"`
}

func PostGetOne(w http.ResponseWriter, r *http.Request) {
	var sid string = mux.Vars(r)["id"]
	id, _ := strconv.Atoi(sid)
	p := md.Post{Id:id}
	if err := p.PostGetOneSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}
	resp :=SinglePostResponse{Sp:p}
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

