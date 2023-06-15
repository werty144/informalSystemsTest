package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		tokens := strings.Fields(scanner.Text())
		if len(tokens) == 0 {
			fmt.Println("Bad input")
			continue
		}
		var command = tokens[0]

		switch command {
		case "stop":
			stop()
		case "start":
			start()
		case "play":
			play()
		default:
			fmt.Println("Unknown command")
			continue
		}
	}
}

func stop() {
	os.Exit(0)
}

func play() { fmt.Println("Playing!") }

func start() { fmt.Println("Starting!") }
