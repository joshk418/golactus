package golactus

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	ln             net.Listener
	name           string
	host           string
	port           string
	server         *http.Server
	router         *mux.Router
	requestTimeout int

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

type MiddlwareHandler func(h http.Handler) http.Handler

func (s *Server) RegisterMiddleware(middleware ...MiddlwareHandler) {
	for _, m := range middleware {
		s.router.Use(mux.MiddlewareFunc(m))
	}
}

func (s *Server) AddRoutes(reqs ...Request) *Server {
	for _, r := range reqs {
		s.router.HandleFunc(r.url, handleFunc(r.handler)).Methods(r.action)
	}

	return s
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
			log.Println(err)
			handleError(err, w)
			return
		}

		log.Infof("[%s] %s", r.Method, r.RequestURI)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		json.NewEncoder(w).Encode(i)
	})
}
