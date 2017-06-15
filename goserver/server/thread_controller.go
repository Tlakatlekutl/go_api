package server

import (
	md "./models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	//"fmt"
	"strconv"
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
		t.Created = "1970-01-01T00:00:00.000Z"
	}

	if err := t.ThreadCreateSQL(DB.DB); err != nil {
		switch err {
		case md.UniqueError:
			if err := t.ThreadGetOneSQL(DB.DB); err == nil {
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
	fu := md.ForumUser{f.Slug, t.Author}

	if  err := fu.ForumUserInsertSQL(DB.DB); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
	}

	RespondJSON(w, http.StatusCreated, t)
}

func ThreadGetOne(w http.ResponseWriter, r *http.Request) {
	var ppk string = mux.Vars(r)["slug_or_id"]
	t := md.Thread{}

	if id, err := IsId(ppk); err == nil {
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

type PostsResponse struct {
	Marker string    `json:"marker, omitempty"`
	Posts  []md.Post `json:"posts"`
}

func ThreadGetPosts(w http.ResponseWriter, r *http.Request) {
	var ppk string = mux.Vars(r)["slug_or_id"]
	t := md.Thread{}

	if id, err := IsId(ppk); err == nil {
		t.ID = id
	} else {
		t.Slug = ppk
	}
	if err := t.ThreadSelectOneIdOrSlugSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	queryParams := r.URL.Query()
	var limit string
	if val, ok := queryParams["limit"]; ok {
		limit = val[0]
	}
	var marker int
	if val, ok := queryParams["marker"]; ok {
		marker, _ = strconv.Atoi(val[0])
	}
	var sort string
	if val, ok := queryParams["sort"]; ok {
		sort = val[0]
	}
	var desc string
	if val, ok := queryParams["desc"]; ok {
		desc = val[0]
	}

	if sort == "flat" || sort == "" {
		if posts, err := t.ThreadGetPostsFlatSQL(DB.DB, limit, desc, marker); err == nil {
			if len(posts) != 0 {
				l, _ := strconv.Atoi(limit)
				marker += l
			}
			resp := PostsResponse{Marker: strconv.Itoa(marker), Posts: posts}
			RespondJSON(w, http.StatusOK, resp)
			return
		} else {
			CheckDbErr(err, w)
			return
		}
	} else if sort == "tree" {
		if posts, err := t.ThreadGetPostsTreeSQL(DB.DB, limit, desc, marker); err == nil {
			if len(posts) != 0 {
				l, _ := strconv.Atoi(limit)
				marker += l
			}
			resp := PostsResponse{Marker: strconv.Itoa(marker), Posts: posts}
			RespondJSON(w, http.StatusOK, resp)
			return
		} else {
			CheckDbErr(err, w)
			return
		}
	} else if sort == "parent_tree" {
		if posts, err := t.ThreadGetPostsParentTreeSQL(DB.DB, limit, desc, marker); err == nil {
			if len(posts) != 0 {
				l, _ := strconv.Atoi(limit)
				marker += l
			}
			resp := PostsResponse{Marker: strconv.Itoa(marker), Posts: posts}
			RespondJSON(w, http.StatusOK, resp)
			return
		} else {
			CheckDbErr(err, w)
			return
		}
	}

}

func ThreadUpdate(w http.ResponseWriter, r *http.Request) {
	var ppk string = mux.Vars(r)["slug_or_id"]
	t := md.Thread{}

	if id, err := IsId(ppk); err == nil {
		t.ID = id
	} else {
		t.Slug = ppk
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid parse thread update json")
		return
	}
	defer r.Body.Close()
	if t.Message == "" && t.Title == "" {
		if err := t.ThreadSelectOneIdOrSlugSQL(DB.DB); err != nil {
			CheckDbErr(err, w)
			return
		}
		RespondJSON(w, http.StatusOK, t)
		return
	}

	if err := t.ThreadUpdateSQL(DB.DB); err != nil {
		CheckDbErr(err, w)
		return
	}

	RespondJSON(w, http.StatusOK, t)

}

func ThreadVote(w http.ResponseWriter, r *http.Request) {
	var ppk string = mux.Vars(r)["slug_or_id"]
	t := md.Thread{}

	if id, err := IsId(ppk); err == nil {
		t.ID = id
	} else {
		t.Slug = ppk
		if err := t.ThreadSelectOneIdOrSlugSQL(DB.DB); err != nil {
			CheckDbErr(err, w)
			return
		}
	}

	v := md.Vote{Thread: t.ID}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid parse vote json")
		return
	}
	defer r.Body.Close()


	if err := t.ThreadVote(DB.DB, &v); err != nil {
		CheckDbErr(err, w)
		return
	}

	RespondJSON(w, http.StatusOK, t)

}
