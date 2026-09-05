// Main entry point for the go-rcon CLI tool
package main

import (
	"fmt"
	"os"

	"github.com/discohaus/go-rcon/internal/cli"
	"github.com/spf13/cobra"
)

func main() {
	var host string
	var port int32
	var password string
	var charset string

	rootCmd := &cobra.Command{
		Use:   "go-rcon",
		Short: "Simple CLI Tool to connect to a RCON server and send Commands. Made by DiscoHaus. See more at https://github.com/discohaus/go-rcon",
		Run: func(_ *cobra.Command, _ []string) {
			cli, err := cli.NewCli(&host, &port, &password, &charset)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "go-rcon: "+err.Error())
				os.Exit(1)
			}
			if err := cli.Run(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "go-rcon: "+err.Error())
				os.Exit(1)
			}
		},
	}
	rootCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Server Host")
	rootCmd.Flags().Int32VarP(&port, "port", "P", 25575, "Server Port which RCON is listening on")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "RCON Password")
	rootCmd.Flags().StringVarP(&charset, "charset", "c", "latin1", "Charset to use for RCON Payloads Options: latin1, utf8, ascii")

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "go-rcon: "+err.Error())
		os.Exit(1)
	}

}
