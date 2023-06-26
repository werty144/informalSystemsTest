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
	/*
		The application entry point. Processes user input in the for loop accepting following commands:
			(*) start --value v --max-value max --num-agents number --liar-ratio ratio
			Launches a number of independent agents, with number * (1-ratio) honest agents always responding
			with the specified integer value v, and (number * ratio) liar agents responding x with x != v and 1 <= x
			<= max. This command starts the game and, when ready, produces the agents.config file, which will contain
			sufficient information to identify and communicate with all agents reliably, and print “ready” on the
			terminal.
			Can be invoked multiple times. When invoked consecutively, cleans previously used resources.

			(*) play
			May be invoked multiple times. Upon each invocation, the client reads the agents.config file (which is
			treated as coming from an external service), connects to the agents using the TCP-based protocol, plays a
			round of the game, and prints the network value v.

			(*) stop
			Stops all the agents listed in the file agents.config, removes the information about stopped agents
			from the file, and exits from the executable
	*/
	rand.Seed(time.Now().UnixNano())
	stopChanel := make(chan struct{}) // chanel used to terminate agent routines

	// process the input
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
			os.Exit(0)
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
	/*
		Stops agents and cleans the agent.config file.
		Arguments:
			stopChanel: the pointer to the [non nil] chanel used to stop agent processes
		Can be called even if there is no running session.
	*/
	close(*stopChanel)          // send signal for agents to quit
	time.Sleep(1 * time.Second) // wait for all agents to gracefully quit
	cleanFile()                 // clean the config file
}

func play() {
	/*
		Plays one round of the game.
	*/
	playGame()
}

func start(networkValue int, maxValue int, nAgents int, liarsRatio float64, stopChanel *chan struct{}) {
	/*
		Launches a number of independent agents, with nAgents * (1-liarsRatio) honest agents always responding
		with the specified integer value networkValue, and (nAgents * liarsRatio) liar agents responding x with
		x != networkValue and 1 <= x <= maxV. This command starts the game and, when ready, produces the agents.config
		file, which will contain sufficient information to identify and communicate with all agents reliably.
		Arguments:
			networkValue: value of correct processes
			maxValue: maximal value for liars
			nAgents: number of agents in the system
			liarsRatio: ratio of liars in the system
			stopChanel: pointer to chanel to stop agent processes
	*/

	stop(stopChanel) // clean resources in case it is not the first call

	*stopChanel = make(chan struct{}) // create a fresh stopping chanel

	// calculate the number of correct and liar agents
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
	// Used here, so it is impossible to retrieve the information about which processes are liars
	// from the order they appear in the config file
	for i := len(addresses) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		addresses[i], addresses[j] = addresses[j], addresses[i]
	}

	createConfigFile(addresses) // create agents.config file
}

func createConfigFile(addresses []*net.TCPAddr) {
	/*
		Creates an agent.config file in the same directory with an executable.
		The file contains addresses in the format <IP>:<port> for each agent, one per line
		Arguments:
			addresses: TCP addresses of agents in the system
	*/
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
	/*
		Parses the start command.
		Arguments:
			tokens: space separated tokens passed by the user in a start command
		Returns:
			err: error in case the input is incorrect
			v: value passed by the user
			maxV: max value passed by the user
			nAgents: number of agents passed by the user
			ratio: liars ratio passed by the user
		Requires strict match with the expected string format
	*/
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
	/*
		Empties the agents.config file
	*/
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
