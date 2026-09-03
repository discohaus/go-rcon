package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/jltobler/go-rcon"
)

func main() {
	host := "localhost"
	port := 25579
	password := "1234"

	client := rcon.NewClient(fmt.Sprintf("rcon://%s:%d", host, port), password, rcon.WithOptions(rcon.CharSetUTF8))

	fmt.Printf("Verbunden mit %s:%d!\n", host, port)
	fmt.Println("Gib ein Kommando ein (tippe 'exit' zum Beenden):")
	fmt.Println("-----------------------------------------------------------------------")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		if scanner.Err() != nil {
			fmt.Printf("Fehler beim Lesen der Eingabe: %v\n", scanner.Err())
			continue
		}

		input := scanner.Text()
		if input == "exit" || input == "quit" {
			fmt.Println("CLI beendet.")
			break
		}

		// Predictions für die Eingabe abrufen

		start := time.Now()
		resp, err := client.Send(input)
		duration := time.Since(start)
		fmt.Printf("Dauer: %v\n", duration)
		if err != nil {
			fmt.Printf("Fehler beim Abrufen der Predictions: %v\n\n", err)
			continue
		}

		fmt.Print(resp)
		fmt.Println()
	}
}
