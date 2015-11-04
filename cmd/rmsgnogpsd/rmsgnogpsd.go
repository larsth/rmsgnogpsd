package main

import (
	"log"
	"os"

	"github.com/larsth/rmsgapps/cmd/rmsgnogpsd/internal/command"
	"github.com/spf13/cobra"
)

func main() {
	var cmd *cobra.Command = command.RootCmd
	var err error

	//log.Logger settings
	log.SetFlags(log.Ldate | log.Lshortfile | log.LUTC)
	log.SetOutput(os.Stderr)
	log.SetPrefix(command.CommandName)

	if cmd == nil {
		log.Fatalln("cmd er <nil>")
	}
	err = cmd.Execute()
	if err != nil {
		os.Exit(-2)
	}
	os.Exit(0)
}
