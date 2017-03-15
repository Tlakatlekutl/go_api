package server

import (
	"net/http"
	//"github.com/gorilla/mux"
	//"encoding/json"
)



//func Clear(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//
//}
//
//func ForumCreate(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//}
//
//func ForumGetOne(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//}
//
//func ForumGetThreads(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//}

func ForumGetUsers(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
}

//func PostGetOne(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//	        w.Write([]byte("hello"))
//		w.Write([]byte(mux.Vars(r)["id"]))
//}
//
//func PostUpdate(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//}
//
//func PostsCreate(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//}

func Status(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
}

//func ThreadCreate(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//}
//
//func ThreadGetOne(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//}
//
//func ThreadGetPosts(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//}
//
//func ThreadUpdate(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//}
//
//func ThreadVote(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//}

//func UserCreate(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//}
//
//func UserGetOne(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.WriteHeader(http.StatusOK)
//	        u := User{mux.Vars(r)["nickname"], "Captain Jack Sparrow",
//			"This is the day you will always remember as the day that you almost caught Captain Jack Sparrow!",
//			"captaina@blackpearl.sea"}
//	        json.NewEncoder(w).Encode(u)
//}
//
//func UserUpdate(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//}

