package main

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/kettek/pingmon/pkg/core"

	"github.com/getlantern/systray"
)

var c *core.Config
var entries []*systray.MenuItem
var lock sync.Mutex

func main() {
	c = &core.DefaultConfig

	if err := c.FromYAML("cfg.yml"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if c.Targets == nil || len(*c.Targets) == 0 {
		fmt.Println(errors.New("no targets"))
		os.Exit(1)
	}

	systray.Run(onReady, onExit)
}

func refreshEntries() {
	lock.Lock()
	defer lock.Unlock()
	for i, t := range *c.Targets {
		title := fmt.Sprintf("%s ", t.Status)
		if t.Name != "" {
			title += t.Name
		} else {
			if *c.ShowMethods {
				title += fmt.Sprintf("%s:", t.Method)
			}
			title += t.Address
			if *c.ShowPorts {
				title += fmt.Sprintf(":%d", t.Port)
			}
		}
		if t.Status == "online" {
			title += fmt.Sprintf(" @ %.3fms", t.Delay/1024)
			entries[i].Check()
		} else {
			entries[i].Uncheck()
		}
		entries[i].SetTitle(title)
		entries[i].SetTooltip(fmt.Sprintf("%s: %s", t.Status, t.ExtendedStatus))
	}
}

func onReady() {
	var prefix, name, suffix string
	if c.Title.Prefix != nil {
		prefix = *c.Title.Prefix
	}
	if c.Title.Name != nil {
		name = *c.Title.Name
	}
	if c.Title.Suffix != nil {
		suffix = *c.Title.Suffix
	}

	systray.SetTitle(fmt.Sprintf("%s%s%s", prefix, name, suffix))

	// Add entries for all our targets.
	for range *c.Targets {
		entries = append(entries, systray.AddMenuItem("", ""))
	}
	systray.AddSeparator()
	// Add our normal menu items.
	refresh := systray.AddMenuItem("Refresh", "Refreshes all targets")
	go func() {
		for {
			<-refresh.ClickedCh
			core.RefreshTargets(c)
			refreshEntries()
		}
	}()
	quit := systray.AddMenuItem("Quit", "Quit")
	go func() {
		for {
			<-quit.ClickedCh
			systray.Quit()
		}
	}()

	// Update and sync.
	core.RefreshTargets(c)
	refreshEntries()

	// Start our refresh timer.
	go func() {
		for {
			core.RefreshTargets(c)
			refreshEntries()
			select {
			case <-time.After(time.Duration(c.Rate) * time.Second):
			}
		}
	}()
}

func onExit() {
	// dunno
}
