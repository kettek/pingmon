package main

import (
	. "github.com/kettek/gobl"
)

func main() {
	// Backend
	Task("buildBackend").
		Exec("go", "build", "-v", "./cmd/pingmon")
	Task("runBackend").
		Exec("./pingmon")
	Task("watchBackend").
		Watch("cmd/pingmon/*", "pkg/backend/*.go", "pkg/core/*.go").
		Signaler(SigQuit).
		Run("buildBackend").
		Run("runBackend")

	// Frontend
	Task("buildFrontend").
		Chdir("pkg/frontend").
		Exec("npm", "i").
		Exec("npm", "run", "build")
	Task("watchFrontend").
		Watch("pkg/frontend/src/*.*").
		Signaler(SigQuit).
		Run("buildFrontend")

	// CLI
	Task("buildCLI").
		Exec("go", "build", "-v", "./cmd/pingmon-cli")
	Task("runCLI").
		Exec("./pingmon-cli")
	Task("watchCLI").
		Watch("cmd/pingmon-cli/*", "pkg/core/*").
		Signaler(SigQuit).
		Run("buildCLI").
		Run("runCLI")

	// Installers
	Task("installSystemdUnit").
		Exec("mkdir", "-p", "/opt/pingmon/public").
		Exec("cp", "pingmon", "/opt/pingmon/").
		Exec("cp", "-r", "pkg/frontend/public", "/opt/pingmon/").
		Exec("cp", "-n", "extra/cfg.yml", "/opt/pingmon/").
		Exec("cp", "extra/pingmon.service", "/etc/systemd/system/")

	Go()
}
