package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/websocket"

	"github.com/go-ping/ping"
)

type Server struct {
	clients       []*wsClient
	clientsLock   sync.Mutex
	server        *http.Server
	lastTimestamp time.Time
	ExitChan      chan struct{}
	running       bool
}

type wsClient struct {
	conn *websocket.Conn
	exit chan struct{}
}

func (s *Server) Run(c *Config) error {
	s.ExitChan = make(chan struct{})

	mux := http.NewServeMux()

	apiServerHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		b, err := json.Marshal(c.Title)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, string(b))
	}
	apiServicesHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if !s.running {
			s.RefreshTargets(c)
		}
		b, err := json.Marshal(struct {
			Elapsed time.Duration `json:"elapsed"`
			Targets *[]*Target    `json:"targets"`
		}{
			Targets: c.Targets,
			Elapsed: time.Duration(time.Now().Sub(s.lastTimestamp).Milliseconds()),
		})
		if err != nil {
			panic(err)
		}
		io.WriteString(w, string(b))
	}
	fileHandler := http.FileServer(http.Dir(*c.Assets))
	wsHandler := websocket.Handler(func(conn *websocket.Conn) {
		cl := &wsClient{
			conn: conn,
			exit: make(chan struct{}),
		}
		if !s.running {
			s.RefreshTargets(c)
		}
		s.UpdateClient(c, cl)
		s.clientsLock.Lock()
		s.clients = append(s.clients, cl)
		s.UpdateClients(c)
		s.clientsLock.Unlock()
		select {
		case <-cl.exit:
		}
	})

	mux.HandleFunc("/api/title", apiServerHandler)
	mux.HandleFunc("/api/services", apiServicesHandler)
	mux.Handle("/ws", wsHandler)
	mux.Handle("/", fileHandler)

	return http.ListenAndServe(*c.Address, mux)
}

func (s *Server) UpdateClients(c *Config) {
	if len(s.clients) == 0 && s.running {
		s.StopLoop()
	} else if !s.running {
		s.StartLoop(c)
	}
}

func (s *Server) StopLoop() {
	go func() {
		s.ExitChan <- struct{}{}
		s.running = false
	}()
}

func (s *Server) StartLoop(c *Config) {
	s.running = true
	go s.PokeLoop(c)
}

func (s *Server) RefreshTargets(c *Config) {
	if c.Targets == nil {
		return
	}
	for _, t := range *c.Targets {
		if t.Method == "tcp" || t.Method == "udp" {
			start := time.Now()
			_, err := net.DialTimeout(t.Method, fmt.Sprintf("%s:%d", t.Address, t.Port), time.Duration(*c.Timeout)*time.Second)
			if err != nil {
				t.Status = "offline"
			} else {
				t.Status = "online"
			}
			t.Delay = float64(time.Now().Sub(start).Microseconds())
		} else if t.Method == "ping" {
			pinger, err := ping.NewPinger(t.Address)
			pinger.SetPrivileged(*c.PrivilegedPing)
			if err != nil {
				t.Status = "error"
				continue
			}
			pinger.Timeout = time.Duration(*c.Timeout) * time.Second
			pinger.Count = 1
			pinger.Interval = 1
			start := time.Now()
			if err := pinger.Run(); err != nil {
				fmt.Println(err)
				t.Status = "offline"
			} else if pinger.Statistics().PacketsRecv == 0 {
				t.Status = "offline"
			} else {
				t.Status = "online"
			}
			t.Delay = float64(time.Now().Sub(start).Microseconds())
		}
	}
	s.lastTimestamp = time.Now()
}

func (s *Server) UpdateClient(c *Config, cl *wsClient) {
	if err := websocket.JSON.Send(cl.conn, struct {
		Elapsed time.Duration `json:"elapsed"`
		Targets *[]*Target    `json:"targets"`
	}{
		Targets: c.Targets,
		Elapsed: time.Duration(time.Now().Sub(s.lastTimestamp).Milliseconds()),
	}); err != nil {
		s.clientsLock.Lock()
		for i, c := range s.clients {
			if c == cl {
				s.clients = append(s.clients[:i], s.clients[i+1:]...)
				break
			}
		}
		s.UpdateClients(c)
		s.clientsLock.Unlock()
		cl.conn.Close()
		cl.exit <- struct{}{}
	}
}

func (s *Server) PokeLoop(c *Config) {
	for {
		// Loop through our defined targets.
		s.RefreshTargets(c)
		// Update our clients.
		for _, cl := range s.clients {
			s.UpdateClient(c, cl)
		}
		// Check if we should bail, otherwise sleep before looping.
		select {
		case <-s.ExitChan:
			return
		case <-time.After(time.Duration(c.Rate) * time.Second):
		}
	}
}
