package main

import (
	. "github.com/kettek/gobl"
)

func main() {
	Task("buildBackend").
		Exec("go", "build", "-v", "./cmd/server")

	Task("runBackend").
		Exec("./server")

	Task("buildFrontend").
		Chdir("pkg/frontend").
		Exec("npm", "run", "build")

	Task("watchBackend").
		Watch("cmd/server/*", "pkg/backend/*.go").
		Signaler(SigQuit).
		Run("buildBackend").
		Run("runBackend")

	Task("watchFrontend").
		Watch("pkg/frontend/src/*.*").
		Signaler(SigQuit).
		Run("buildFrontend")

	Go()
}
