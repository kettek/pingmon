package main

import (
	"log"
	"pingmon/pkg/backend"
	"pingmon/pkg/cfg"
)

func main() {
	s := backend.Server{}
	c := &cfg.Config{
		Address: ":8999",
	}

	if err := c.FromYAML("cfg.yml"); err != nil {
		log.Println(err)
		log.Println("using defaults")
	}

	if err := s.Run(c); err != nil {
		log.Fatal(err)
	}
	log.Println("running...")
}
