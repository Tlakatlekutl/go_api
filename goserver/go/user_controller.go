package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	md "./models"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	var nn string = mux.Vars(r)["nickname"]
	u := md.User{Nickname:nn}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request creation user")
		return
	}
	defer r.Body.Close()

	if err := u.CreateUserSQL(DB.DB); err != nil {
		switch err {
			case md.UniqueError:
				if users, err := u.GetUniqueUsersSQL(DB.DB); err == nil {
					RespondJSON(w, http.StatusConflict, users)
					return
				}

			default:
				RespondError(w, http.StatusInternalServerError, err.Error())
			}
		return
	}

	RespondJSON(w, http.StatusCreated, u)

}

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	var nn string = mux.Vars(r)["nickname"]
	u := md.User{Nickname:nn}

	if err := u.GetOneUserSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	RespondJSON(w, http.StatusOK, u)
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	var nn string = mux.Vars(r)["nickname"]
	u := md.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request update user")
		return
	}
	defer r.Body.Close()

	u.Nickname = nn
	 if InspectEmpty(u) {
		 u2 := md.User{Nickname:nn}
		 if err := u2.GetOneUserSQL(DB.DB); err != nil {
			 CheckDbErr(err, w)
			 return
		 }
		 CompareTypes(&u, &u2)
	 }

	if err := u.UpdateUserSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	RespondJSON(w, http.StatusOK, u)
}
