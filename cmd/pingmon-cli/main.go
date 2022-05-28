package main

import (
	"errors"
	"fmt"
	"os"
	"pingmon/pkg/core"
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
		fmt.Printf("%s: %s @ %fms\n", t.Address, t.Status, t.Delay/1024)
		if t.ExtendedStatus != "" {
			fmt.Printf("\t%s\n", t.ExtendedStatus)
		}
	}
}
