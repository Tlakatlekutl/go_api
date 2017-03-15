package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	md "./models"
	//"fmt"
)

func PostGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
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

