package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var stopChanel chan struct{}

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
			stop(&stopChanel)
		case "start":
			start(&stopChanel)
		case "play":
			play()
		default:
			fmt.Println("Unknown command")
			continue
		}
	}
}

func stop(stopChanel *chan struct{}) {
	close(*stopChanel)
}

func play() { fmt.Println("Playing!") }

func start(stopChanel *chan struct{}) {
	*stopChanel = make(chan struct{})
	fmt.Println("Starting!")
	go echo(*stopChanel)
}
