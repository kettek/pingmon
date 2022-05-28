package backend

import (
	"time"

	"github.com/kettek/pingmon/pkg/core"

	"golang.org/x/net/websocket"
)

type wsClient struct {
	conn *websocket.Conn
	exit chan struct{}
}

func (s *Server) HandleWS(conn *websocket.Conn) {
	cl := &wsClient{
		conn: conn,
		exit: make(chan struct{}),
	}
	if !s.running {
		s.RefreshTargets()
	}
	s.updateClient(cl)
	s.clientsLock.Lock()
	s.clients = append(s.clients, cl)
	s.checkLoop()
	s.clientsLock.Unlock()
	select {
	case <-cl.exit:
	}
}

func (s *Server) updateClient(cl *wsClient) {
	if err := websocket.JSON.Send(cl.conn, struct {
		Elapsed time.Duration   `json:"elapsed"`
		Targets *[]*core.Target `json:"targets"`
	}{
		Targets: s.config.Targets,
		Elapsed: time.Duration(time.Now().Sub(s.lastTimestamp).Milliseconds()),
	}); err != nil {
		s.clientsLock.Lock()
		for i, c := range s.clients {
			if c == cl {
				s.clients = append(s.clients[:i], s.clients[i+1:]...)
				break
			}
		}
		s.checkLoop()
		s.clientsLock.Unlock()
		cl.conn.Close()
		cl.exit <- struct{}{}
	}
}
