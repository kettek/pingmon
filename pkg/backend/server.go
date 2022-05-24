package backend

import (
	"io"
	"net/http"
	"pingmon/pkg/cfg"
)

type Server struct {
	server *http.Server
}

func (s *Server) Run(c *cfg.Config) error {
	apiHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Yo\n")
	}
	fileHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Yo file\n")
	}

	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/", fileHandler)

	return http.ListenAndServe(c.Address, nil)
}
