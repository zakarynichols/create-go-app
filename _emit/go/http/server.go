package http

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	thing "github.com/username/repo"
	"github.com/username/repo/cors"
)

type Server struct {
	ln           net.Listener
	server       *http.Server
	router       *mux.Router
	thingService thing.ThingService
}

type Config struct {
	Addr         string
	ThingService thing.ThingService
}

func New(config Config) *Server {
	r := mux.NewRouter()
	c := cors.New()
	h := c.Handler(r)
	return &Server{
		server: &http.Server{
			Addr:    ":" + config.Addr,
			Handler: h,
		},
		router:       r,
		thingService: config.ThingService,
	}
}

func (s *Server) Open() error {
	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	s.ln = ln

	err = s.server.Serve(s.ln)
	if err != nil {
		return err
	}

	return nil
}
