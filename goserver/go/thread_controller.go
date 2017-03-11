package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	md "./models"
)


func ThreadCreate(w http.ResponseWriter, r *http.Request) {
	var fsg string = mux.Vars(r)["slug"]
	t := md.Thread{Forum: fsg}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request creation thread")
		return
	}
	defer r.Body.Close()

	if (t.Created == "") {
		t.Created="2017-08-22T01:30:51.934+03:00"
	}
	if err := t.ThreadCreateSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	RespondJSON(w, http.StatusCreated, t)
}

func ThreadGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadGetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ThreadVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
