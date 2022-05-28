package main

import (
	"log"
	"pingmon/pkg/backend"
	"pingmon/pkg/core"
)

func main() {
	s := backend.Server{}
	c := &core.DefaultConfig

	if err := c.FromYAML("cfg.yml"); err != nil {
		log.Println(err)
		log.Println("using defaults")
	}

	log.Printf("serving files in \"%s\" on \"%s\"\n", *c.Assets, *c.Address)

	if c.Targets != nil {
		for _, t := range *c.Targets {
			log.Println(t)
		}
	}

	if err := s.Run(c); err != nil {
		log.Fatal(err)
	}
}
