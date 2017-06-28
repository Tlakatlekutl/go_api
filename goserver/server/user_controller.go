package server

import (
	md "./models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
)

func (a *App) UserCreate(w http.ResponseWriter, r *http.Request) {
	var nn string = mux.Vars(r)["nickname"]
	u := md.User{Nickname: nn}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request creation user")
		return
	}
	defer r.Body.Close()

	if err := u.CreateUserSQL(a.DB); err != nil {
		switch err {
		case md.UniqueError:
			if users, err := u.GetUniqueUsersSQL(a.DB); err == nil {
				RespondJSON(w, http.StatusConflict, users)
				return
			} else {
				fmt.Println(err.Error())
				RespondError(w, http.StatusRequestTimeout, err.Error())
			}

		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondJSON(w, http.StatusCreated, u)

}

func (a *App) UserGetOne(w http.ResponseWriter, r *http.Request) {
	var nn string = mux.Vars(r)["nickname"]
	u := md.User{Nickname: nn}

	if err := u.GetOneUserSQL(a.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	RespondJSON(w, http.StatusOK, u)
}

func (a *App) UserUpdate(w http.ResponseWriter, r *http.Request) {
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
		u2 := md.User{Nickname: nn}
		if err := u2.GetOneUserSQL(a.DB); err != nil {
			CheckDbErr(err, w)
			return
		}
		CompareTypes(&u, &u2)
	}

	if err := u.UpdateUserSQL(a.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	RespondJSON(w, http.StatusOK, u)
}
