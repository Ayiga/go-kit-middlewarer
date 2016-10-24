package gorilla

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router is a type alias for *github.com/gorilla/mux.Router to make reference
// easier, as well as to implement the Mux interface.
type Router struct {
	*mux.Router
}

type CastError int

// implements the error interface.
func (CastError) Error() string {
	return "Unable to cast Router to *mux.Router"
}

// NewRouter returns a new Router, wrapping the given gorilla mux Router.
func NewRouter(router *mux.Router) *Router {
	return &Router{
		Router: router,
	}
}

// Handle registers the handler for the given pattern.
// According to net/http.ServeMux If a handler already exists for pattern,
// the Handle invocation panics.
func (r *Router) Handle(pattern string, handler http.Handler) {
	r.Router.Handle(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern.
func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.Router.HandleFunc(pattern, handler)
}

// ServerHTTP implements net/http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}
