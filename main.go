package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
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
			time.Sleep(1 * time.Second)
			log.Println("Stopped.")
			//os.Exit(0)
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
	go startServer(*stopChanel)
}
