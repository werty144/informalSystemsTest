package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func playGame() {
	/*
		Plays a round of the game.
		Connects to the agents using addresses from the agent.config file,
		collects their values and prints the mode value among collected.

		In case liars ratio in the system is >= 0.5 the printed value is NOT guaranteed to be the network value.
	*/

	addresses := getAgentAddresses()
	err, networkValue := getNetworkValue(addresses)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("The network value is %d\n", networkValue)
}

func getNetworkValue(addresses []net.TCPAddr) (error, int) {
	/*
		Queries s the agents for their values and returns the mode of the collected values.
		Arguments:
			addresses: TCP addresses of the agents in the system
		Returns:
			err: Error in case failed to collect all the values
			network value: network value
		In case liars ratio in the system is >= 0.5 the returned value is NOT guaranteed to be the network value.
	*/

	var values ThreadSafeSlice // thread safe values storage

	// used to asynchronously query the agents
	var wg sync.WaitGroup
	wg.Add(len(addresses))

	for _, addr := range addresses {
		go getAgentsValue(addr, &values, &wg)
	}

	wg.Wait()

	if len(values.slice) != len(addresses) {
		// failed to collect all values
		return errors.New(fmt.Sprintf("Received %d out of %d values", len(values.slice), len(addresses))), 0
	}

	networkValue := getMode(values.slice) // get the mode of collected values
	return nil, networkValue
}

func getAgentsValue(addr net.TCPAddr, values *ThreadSafeSlice, wg *sync.WaitGroup) {
	/*
		Query the agent for its value.
		In case of successfully obtaining agents value, adds it to the values storage.
		Arguments:
			addr: TCP address of the agent
			values: shared thread safe values storage
			wg: wait group for signalizing the termination
		Fault tolerant: if there is no response within a second returns adding nothing to the values storage.
	*/

	defer wg.Done() // when done signalize the termination

	// connect to the agent with timeout of 1 sec
	conn, err := net.DialTimeout("tcp", addr.String(), 1*time.Second)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer conn.Close()

	// send an arbitrary message to the agent
	_, err = conn.Write([]byte("ping"))
	if err != nil {
		log.Fatal("Error sending message:", err)
	}

	response := make([]byte, 1024) // Assuming the response size is less than 1024 bytes
	length, err := conn.Read(response)
	if err != nil {
		// might be caused because of the connection timeout
		fmt.Println("Error receiving response from "+addr.String()+".", err)
		return
	}

	v := binary.BigEndian.Uint32(response[:length])
	values.Append(int(v)) // thread safely append the value to the common storage
}

func getAgentAddresses() []net.TCPAddr {
	/*
		Collect agent addresses from the agent.config file
		Returns:
			Agents TCP addresses from agent.config file
	*/

	file, err := os.Open("agents.config")
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer file.Close()

	var addresses []net.TCPAddr

	scanner := bufio.NewScanner(file)
	// Collect addresses from the file assuming there is one address per line of the form <IP>:<port>
	for scanner.Scan() {
		line := scanner.Text()
		tcpAddr, err := net.ResolveTCPAddr("tcp", line)
		if err != nil {
			log.Fatal("Error:", err)
		}
		addresses = append(addresses, *tcpAddr)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error:", err)
	}
	return addresses
}

type ThreadSafeSlice struct {
	/*
		Structure used to concurrently add integers.
		Mutex-based.
	*/
	slice []int      // actual data
	mutex sync.Mutex // mutex used for managing shared access
}

func (t *ThreadSafeSlice) Append(value int) {
	/*
		Appends a value to the storage.
		Arguments:
			value: value to be appended
		NOT lock-free
	*/

	t.mutex.Lock() // acquire mutex
	defer t.mutex.Unlock()

	t.slice = append(t.slice, value)
}

func getMode(values []int) int {
	/*
		Get the most frequent element (the mode) of the collection
		Arguments:
			values: slice of the elements among which to search for mode
		Returns:
			the mode of the values
		In case there are several modes, an arbitrary one is returned
	*/

	frequency := make(map[int]int) // introduce a map to count frequencies of the elements
	maxFrequency := 0              // currently maximal frequency
	var mostFrequent int           // current mode

	// Count the frequency of each element
	// maintaining maximal frequency and the current mode
	for _, value := range values {
		frequency[value]++
		if frequency[value] > maxFrequency {
			maxFrequency = frequency[value]
			mostFrequent = value
		}
	}

	return mostFrequent
}
