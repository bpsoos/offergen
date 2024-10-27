package main

import (
	"offergen/cmd"
	"offergen/logging"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		logging.GetLogger().Error("invalid params, usage: 'offergen <entrypoint>' (entrypoint: 'serve' or 'migrate')")
	}

	cmd.NewRootCmd().Execute(os.Args[1])
}
