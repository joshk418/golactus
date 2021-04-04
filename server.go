package golactus

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	ln     net.Listener
	server *http.Server
	router *mux.Router
	name   string
	host   string
	port   string

	Addr   string
	Domain string
}

type Request struct {
	action  string
	url     string
	handler CustomHandlerFunc
}

type CustomHandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func NewServer(opts ...Option) *Server {
	s := &Server{}

	for _, opt := range opts {
		opt(s)
	}

	s.router = mux.NewRouter().PathPrefix("/" + s.name).Subrouter()
	s.server = &http.Server{
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      s.router,
		Addr:         fmt.Sprintf("%s:%s", s.host, s.port),
	}

	return s
}

func (s *Server) AddRoutes(reqs ...Request) {
	for _, r := range reqs {
		s.router.HandleFunc(r.url, handleFunc(r.handler)).Methods(r.action)
	}
}

func (s *Server) Serve() error {
	return s.server.ListenAndServe()
}

func Get(url string, handler CustomHandlerFunc) Request {
	return Request{
		action:  "GET",
		url:     url,
		handler: handler,
	}
}

func Post(url string, handler CustomHandlerFunc) Request {
	return Request{
		action:  "POST",
		url:     url,
		handler: handler,
	}
}

func handleFunc(handler CustomHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i, err := handler(w, r)
		if err != nil {
			handleError(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		json.NewEncoder(w).Encode(i)
	})
}
