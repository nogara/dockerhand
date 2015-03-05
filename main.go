package main

import (
	"runtime"

	"github.com/edgard/dockerhand/commands"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	commands.Execute()
}
