// Main entry point for the go-rcon CLI tool
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/discohaus/go-rcon/internal/cli"
	"github.com/spf13/cobra"
)

func main() {
	var host string
	var port int32
	var password string
	var charset string
	var interactive bool
	var command string
	var file string

	rootCmd := &cobra.Command{
		Use:   "go-rcon",
		Short: "Simple CLI Tool to connect to a RCON server and send Commands. Made by DiscoHaus. See more at https://github.com/discohaus/go-rcon",
		Run: func(cmd *cobra.Command, _ []string) {
			if !interactive && strings.TrimSpace(command) == "" {
				_, _ = fmt.Fprintln(os.Stderr, "go-rcon: either --interactive (-i) or --command (-c) is required")
				os.Exit(1)
			}
			if interactive && strings.TrimSpace(command) != "" {
				_, _ = fmt.Fprintln(os.Stderr, "go-rcon: cannot use both --interactive (-i) and --command (-c)")
				os.Exit(1)
			}

			if file != "" {
				cfg, err := cli.ReadRconProperties(file)
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, "go-rcon: "+err.Error())
					os.Exit(1)
				}
				if !cmd.Flags().Changed("password") && cfg.Password != nil {
					password = *cfg.Password
				}
				if !cmd.Flags().Changed("port") && cfg.Port != nil {
					port = *cfg.Port
				}
			}

			cliInstance, err := cli.NewCli(&host, &port, &password, &charset)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "go-rcon: "+err.Error())
				os.Exit(1)
			}

			if interactive {
				if err := cliInstance.Run(); err != nil {
					_, _ = fmt.Fprintln(os.Stderr, "go-rcon: "+err.Error())
					os.Exit(1)
				}
				return
			}

			output, err := cliInstance.ExecuteCommand(command)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "go-rcon: "+err.Error())
				os.Exit(1)
			}
			if output != "" {
				fmt.Println(output)
			}
		},
	}
	rootCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Server Host")
	rootCmd.Flags().Int32VarP(&port, "port", "P", 25575, "Server Port which RCON is listening on")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "RCON Password")
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "Path to properties file containing rcon.password and rcon.port")
	rootCmd.Flags().StringVarP(&charset, "charset", "C", "latin1", "Charset to use for RCON Payloads Options: latin1, utf8, ascii")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Start interactive TUI session")
	rootCmd.Flags().StringVarP(&command, "command", "c", "", "Execute single RCON command and print output")

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "go-rcon: "+err.Error())
		os.Exit(1)
	}

}
