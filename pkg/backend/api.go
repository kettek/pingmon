package backend

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/kettek/pingmon/pkg/core"
)

func (s *Server) handleSettingsAPI(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	b, err := json.Marshal(struct {
		Title       core.Title `json:"title"`
		ShowMethods bool       `json:"showMethods"`
		ShowPorts   bool       `json:"showPorts"`
	}{
		Title:       *s.config.Title,
		ShowMethods: *s.config.ShowMethods,
		ShowPorts:   *s.config.ShowPorts,
	})
	if err != nil {
		panic(err)
	}
	io.WriteString(w, string(b))
}

func (s *Server) handleServicesAPI(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if !s.running {
		s.RefreshTargets()
	}
	b, err := json.Marshal(struct {
		Elapsed time.Duration   `json:"elapsed"`
		Targets *[]*core.Target `json:"targets"`
	}{
		Targets: s.config.Targets,
		Elapsed: time.Duration(time.Now().Sub(s.lastTimestamp).Milliseconds()),
	})
	if err != nil {
		panic(err)
	}
	io.WriteString(w, string(b))
}
