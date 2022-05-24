package main

import (
	. "github.com/kettek/gobl"
)

func main() {
	Task("buildBackend").
		Exec("go", "build", "-v", "./cmd/pingmon")

	Task("runBackend").
		Exec("./pingmon")

	Task("buildFrontend").
		Chdir("pkg/frontend").
		Exec("npm", "i").
		Exec("npm", "run", "build")

	Task("watchBackend").
		Watch("cmd/pingmon/*", "pkg/backend/*.go").
		Signaler(SigQuit).
		Run("buildBackend").
		Run("runBackend")

	Task("watchFrontend").
		Watch("pkg/frontend/src/*.*").
		Signaler(SigQuit).
		Run("buildFrontend")

	Task("installSystemdUnit").
		Exec("mkdir", "-p", "/opt/pingmon/public").
		Exec("cp", "pingmon", "/opt/pingmon/").
		Exec("cp", "-r", "pkg/frontend/public", "/opt/pingmon/").
		Exec("cp", "-n", "extra/cfg.yml", "/opt/pingmon/").
		Exec("cp", "extra/pingmon.service", "/etc/systemd/system/")

	Go()
}
