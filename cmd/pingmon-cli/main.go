package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/kettek/pingmon/pkg/core"
)

func main() {
	c := &core.DefaultConfig

	if err := c.FromYAML("cfg.yml"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if c.Targets == nil || len(*c.Targets) == 0 {
		fmt.Println(errors.New("no targets"))
		os.Exit(1)
	}

	core.RefreshTargets(c)

	for _, t := range *c.Targets {
		if t.Status == "online" {
			fmt.Printf("👌 ")
		} else if t.Status == "offline" {
			fmt.Printf("👎 ")
		} else if t.Status == "error" {
			fmt.Printf("👎 ")
		} else {
			fmt.Printf("✋ ")
		}
		if t.Name != "" {
			fmt.Printf(t.Name)
		} else {
			if *c.ShowMethods {
				fmt.Printf("%s:", t.Method)
			}
			fmt.Printf(t.Address)
			if *c.ShowPorts {
				fmt.Printf(":%d", t.Port)
			}
		}
		fmt.Printf(" @ %fms\n", t.Delay/1024)
		if t.ExtendedStatus != "" {
			fmt.Printf("  👉 %s\n", t.ExtendedStatus)
		}
	}
}
