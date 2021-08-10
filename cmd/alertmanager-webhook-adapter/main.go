package main

import (
	"fmt"
	"os"

	"github.com/bougou/alertmanager-webhook-adapter/cmd/alertmanager-webhook-adapter/app"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}

func main() {
	command := app.NewRootCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
