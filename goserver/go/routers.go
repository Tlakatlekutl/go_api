package server

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"

)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	//fs := http.FileServer(http.Dir("../static"))
	//
	//router.
	//	Methods("GET").
	//	Path("/static").
	//	Name("Static").
	//	Handler(Logger(fs, "lala"))

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
//func Swagger(w http.ResponseWriter, r *http.Request) {
//	//fmt.Fprintf(w, "Swagger")
//	//w.Write(http.File("/home/tlakatlekutl/GoglandProjects/tech-db-forum/swagger.yml"))
//	//http.HandleFunc
//	//w.Header().Set("Access-Control-Allow-Headers:", "Origin, X-Atmosphere-tracking-id, X-Atmosphere-Framework, X-Cache-Date, Content-Type, X-Atmosphere-Transport, *")
//	w.Header().Set("Access-Control-Allow-Origin:", "\"*\"")
//	//w.Header().Set("Access-Control-Allow-Methods:", "POST, GET, OPTIONS , PUT")
//	//w.Header().Set("Access-Control-Request-Headers:", "Origin, X-Atmosphere-tracking-id, X-Atmosphere-Framework, X-Cache-Date, Content-Type, X-Atmosphere-Transport,  *")
//	//
//	//w.Header().Set("Content-Type:", "application/json")
//	//w.Header().Set("accept:","application/json; charset=utf-8,*/*")
//	//http.ServeFile(w, r, "swagger.yml")
//
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusOK)
//	data, _ := ioutil.ReadFile("swagger.yml")
//	w.Write(data)
//}
func Static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
	fmt.Println(r.URL.Path[1:])
	//w.Header().Set("Access-Control-Allow-Origin:", "*")
	//w.Header().Set("Access-Control-Allow-Methods: GET", "POST, DELETE, PUT, PATCH, OPTIONS")
	//w.Header().Set("Access-Control-Allow-Headers:", "Content-Type, api_key, Authorization")
	//http.ServeFile(w, r, "swagger.json")
}
var routes = Routes{
	Route{
		"Index",
		"GET",
		"/api/",
		Index,
	},
	Route{
		"Static",
		"GET",
		"/static/dist/{rest}",
		Static,
	},
	//

	Route{
		"Clear",
		"POST",
		"/api/service/clear",
		Clear,
	},

	Route{
		"ForumCreate",
		"POST",
		"/api/forum/create",
		ForumCreate,
	},

	Route{
		"ForumGetOne",
		"GET",
		"/api/forum/{slug}/details",
		ForumGetOne,
	},

	Route{
		"ForumGetThreads",
		"GET",
		"/api/forum/{slug}/threads",
		ForumGetThreads,
	},

	Route{
		"ForumGetUsers",
		"GET",
		"/api/forum/{slug}/users",
		ForumGetUsers,
	},

	Route{
		"PostGetOne",
		"GET",
		"/api/post/{id}/details",
		PostGetOne,
	},

	Route{
		"PostUpdate",
		"POST",
		"/api/post/{id}/details",
		PostUpdate,
	},

	Route{
		"PostsCreate",
		"POST",
		"/api/thread/{slug_or_id}/create",
		PostsCreate,
	},

	Route{
		"Status",
		"GET",
		"/api/service/status",
		Status,
	},

	Route{
		"ThreadCreate",
		"POST",
		"/api/forum/{slug}/create",
		ThreadCreate,
	},

	Route{
		"ThreadGetOne",
		"GET",
		"/api/thread/{slug_or_id}/details",
		ThreadGetOne,
	},

	Route{
		"ThreadGetPosts",
		"GET",
		"/api/thread/{slug_or_id}/posts",
		ThreadGetPosts,
	},

	Route{
		"ThreadUpdate",
		"POST",
		"/api/thread/{slug_or_id}/details",
		ThreadUpdate,
	},

	Route{
		"ThreadVote",
		"POST",
		"/api/thread/{slug_or_id}/vote",
		ThreadVote,
	},

	Route{
		"UserCreate",
		"POST",
		"/api/user/{nickname}/create",
		UserCreate,
	},

	Route{
		"UserGetOne",
		"GET",
		"/api/user/{nickname}/profile",
		UserGetOne,
	},

	Route{
		"UserUpdate",
		"POST",
		"/api/user/{nickname}/profile",
		UserUpdate,
	},

}