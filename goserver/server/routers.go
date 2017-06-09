package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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
		//handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
		//Schemes("https, http")
	}
	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/api/",
		Index,
	},

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
