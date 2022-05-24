package backend

import (
	"io"
	"net/http"
)

type Server struct {
	server *http.Server
}

func (s *Server) Run(c *Config) error {
	mux := http.NewServeMux()

	apiHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Yo\n")
	}
	fileHandler := http.FileServer(http.Dir(*c.Assets))

	mux.HandleFunc("/api", apiHandler)
	mux.Handle("/", fileHandler)

	return http.ListenAndServe(*c.Address, mux)
}
