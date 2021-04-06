package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bougou/alertmanager-webhook-adapter/cmd/alertmanager-webhook-adapter/app/options"
	"github.com/bougou/alertmanager-webhook-adapter/pkg/api"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/spf13/cobra"
)

func init() {
	// append the user-defined functions in the command's initialization
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}

func main() {
	run()
}

var rootCmd = &cobra.Command{
	Use:   "alertmanager-webhook-adapter",
	Short: "alertmanager-webhook-adapter",
	Long:  `alertmanager-webhook-adapter`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		addr, _ := cmd.Flags().GetString("listen-address")
		signature, _ := cmd.Flags().GetString("signature")

		appOptions := &options.AppOptions{
			Addr:      addr,
			Signature: signature,
		}

		container := restful.DefaultContainer
		controller := api.NewController(appOptions)
		controller.Install(container)

		s := &http.Server{
			Addr:    appOptions.Addr,
			Handler: container,
		}
		log.Printf("start listening, %s", s.Addr)
		log.Fatal(s.ListenAndServe())

	},
}

func run() {
	rootCmd.Flags().StringP("listen-address", "l", "0.0.0.0:8090", "The address to listen")
	rootCmd.Flags().StringP("signature", "s", "未知", "the signature")
	rootCmd.Flags().AddGoFlagSet(flag.CommandLine)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
