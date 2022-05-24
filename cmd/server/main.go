package main

import (
	"log"
	"pingmon/pkg/backend"
	"pingmon/pkg/cfg"
)

func main() {
	c := &cfg.Config{}

	s := backend.Server{}

	if err := s.Run(c); err != nil {
		log.Fatal(err)
	}
	log.Println("running...")
}
