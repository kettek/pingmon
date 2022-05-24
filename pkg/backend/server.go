package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/go-ping/ping"
)

type Server struct {
	server   *http.Server
	ExitChan chan struct{}
}

func (s *Server) Run(c *Config) error {
	s.ExitChan = make(chan struct{})

	mux := http.NewServeMux()

	apiServerHandler := func(w http.ResponseWriter, req *http.Request) {
		b, err := json.Marshal(c.Title)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, string(b))
	}
	apiServicesHandler := func(w http.ResponseWriter, req *http.Request) {
		b, err := json.Marshal(c.Targets)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, string(b))
	}
	fileHandler := http.FileServer(http.Dir(*c.Assets))

	mux.HandleFunc("/api/title", apiServerHandler)
	mux.HandleFunc("/api/services", apiServicesHandler)
	mux.Handle("/", fileHandler)

	go s.PokeLoop(c)

	return http.ListenAndServe(*c.Address, mux)
}

func (s *Server) PokeLoop(c *Config) {
	fmt.Println("starting poke loop")
	for {
		select {
		case <-s.ExitChan:
			return
		default:
		}
		// Loop through our defined targets.
		if c.Targets != nil {
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
		}

		// Sleep a while.
		time.Sleep(time.Duration(c.Rate) * time.Second)
	}
}
