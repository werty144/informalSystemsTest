package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
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
		case "start":
			err, networkValue, maxValue, nAgents, liarsRatio := parseStartCommand(tokens)
			if err != nil {
				log.Println(err)
				continue
			}
			start(networkValue, maxValue, nAgents, liarsRatio, &stopChanel)
			fmt.Println("ready")
		case "play":
			play()
		default:
			log.Println("Unknown command")
			continue
		}
	}
}

func stop(stopChanel *chan struct{}) {
	close(*stopChanel)
	time.Sleep(1 * time.Second)
	cleanFile()
	os.Exit(0)
}

func play() {
	playGame()
}

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

	for _, address := range addresses {
		_, err = file.WriteString(address.String() + "\n")
		if err != nil {
			log.Fatal("Error:", err)
		}
	}
}

func parseStartCommand(tokens []string) (error, int, int, int, float64) {
	usage := "start --value v --max-value max --num-agents number --liar-ratio ratio"
	if len(tokens) != 9 {
		return errors.New("wrong number of tokens. Usage: " + usage), 0, 0, 0, .0
	}
	if tokens[1] != "--value" ||
		tokens[3] != "--max-value" ||
		tokens[5] != "--num-agents" ||
		tokens[7] != "--liar-ratio" {
		return errors.New("wrong keys. Usage: " + usage), 0, 0, 0, .0
	}

	v, err := strconv.Atoi(tokens[2])
	if err != nil {
		return errors.New("v should be an integer. Usage: " + usage), 0, 0, 0, .0
	}
	maxV, err := strconv.Atoi(tokens[4])
	if err != nil {
		return errors.New("max should be an integer. Usage: " + usage), 0, 0, 0, .0
	}
	nAgents, err := strconv.Atoi(tokens[6])
	if err != nil {
		return errors.New("number of agents should be an integer. Usage: " + usage), 0, 0, 0, .0
	}
	ratio, err := strconv.ParseFloat(tokens[8], 64)
	if err != nil {
		return errors.New("ratio should be a float. Usage: " + usage), 0, 0, 0, .0
	}
	return nil, v, maxV, nAgents, ratio
}

func cleanFile() {
	file, err := os.OpenFile("agents.config", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		log.Fatal("Error:", err)
	}
}
