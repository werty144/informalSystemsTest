package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
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
			start(3, 7, 5, 0.4, &stopChanel)
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

func start(networkValue int, maxValue int, nAgents int, liarsRatio float64, stopChanel *chan struct{}) {
	*stopChanel = make(chan struct{})
	nLiars := int(math.Round(float64(nAgents) * liarsRatio))
	nFair := nAgents - nLiars
	var addresses []*net.TCPAddr

	for i := 0; i < nLiars; i++ {
		addresses = append(addresses, startAgent(networkValue, maxValue, true, *stopChanel))
	}
	for i := 0; i < nFair; i++ {
		addresses = append(addresses, startAgent(networkValue, maxValue, false, *stopChanel))
	}

	// Fisher-Yates algorithm for shuffling
	for i := len(addresses) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		addresses[i], addresses[j] = addresses[j], addresses[i]
	}

	createConfigFile(addresses)
}

func createConfigFile(addresses []*net.TCPAddr) {
	file, err := os.OpenFile("agents.config", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		log.Fatal("Error:", err)
	}

	for _, address := range addresses {
		_, err = file.WriteString(address.String() + "\n")
		if err != nil {
			log.Fatal("Error:", err)
		}
	}
}
