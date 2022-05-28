package backend

import (
	"net/http"
	"sync"
	"time"

	"github.com/kettek/pingmon/pkg/core"

	"golang.org/x/net/websocket"
)

type Server struct {
	config        *core.Config
	clients       []*wsClient
	clientsLock   sync.Mutex
	server        *http.Server
	lastTimestamp time.Time
	exitChan      chan struct{}
	running       bool
}

func (s *Server) Run(c *core.Config) error {
	s.config = c
	s.exitChan = make(chan struct{})

	mux := http.NewServeMux()

	fileHandler := http.FileServer(http.Dir(*c.Assets))

	mux.HandleFunc("/api/title", s.handleTitleAPI)
	mux.HandleFunc("/api/services", s.handleServicesAPI)
	mux.Handle("/ws", websocket.Handler(s.HandleWS))
	mux.Handle("/", fileHandler)

	return http.ListenAndServe(*c.Address, mux)
}

func (s *Server) checkLoop() {
	if len(s.clients) == 0 && s.running {
		s.stopLoop()
	} else if !s.running {
		s.startLoop()
	}
}

func (s *Server) stopLoop() {
	go func() {
		s.exitChan <- struct{}{}
		s.running = false
	}()
}

func (s *Server) startLoop() {
	s.running = true
	go s.PokeLoop()
}

func (s *Server) RefreshTargets() {
	if s.config.Targets == nil {
		s.lastTimestamp = time.Now()
		return
	}
	core.RefreshTargets(s.config)
	s.lastTimestamp = time.Now()
}

func (s *Server) PokeLoop() {
	for {
		// Loop through our defined targets.
		s.RefreshTargets()
		// Update our clients.
		for _, cl := range s.clients {
			s.updateClient(cl)
		}
		// Check if we should bail, otherwise sleep before looping.
		select {
		case <-s.exitChan:
			return
		case <-time.After(time.Duration(s.config.Rate) * time.Second):
		}
	}
}
