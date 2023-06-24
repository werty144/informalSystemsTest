package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func playGame() {
	addresses := getAgentAddresses()

	var values ThreadSafeSlice
	var wg sync.WaitGroup
	wg.Add(len(addresses))

	for _, addr := range addresses {
		go getAgentsValue(addr, &values, &wg)
	}
	wg.Wait()
	if len(values.slice) != len(addresses) {
		log.Printf("Received %d out of %d values", len(values.slice), len(addresses))
		return
	}

	networkValue := getMode(values.slice)
	fmt.Printf("The network value is %d\n", networkValue)
}

func getAgentsValue(addr net.TCPAddr, values *ThreadSafeSlice, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := net.DialTimeout("tcp", addr.String(), 1*time.Second)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("ping"))
	if err != nil {
		log.Fatal("Error sending message:", err)
	}

	response := make([]byte, 1024) // Assuming the response size is less than 1024 bytes
	length, err := conn.Read(response)
	if err != nil {
		fmt.Println("Error receiving response from "+addr.String()+".", err)
	}
	v := binary.BigEndian.Uint32(response[:length])
	values.Append(int(v))
}

func getAgentAddresses() []net.TCPAddr {
	file, err := os.Open("agents.config")
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer file.Close()

	var addresses []net.TCPAddr

	scanner := bufio.NewScanner(file)
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
	slice []int
	mutex sync.Mutex
}

func (t *ThreadSafeSlice) Append(value int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.slice = append(t.slice, value)
}

func getMode(values []int) int {
	frequency := make(map[int]int)
	maxFrequency := 0
	var mostFrequent int

	// Count the frequency of each element
	for _, value := range values {
		frequency[value]++
		if frequency[value] > maxFrequency {
			maxFrequency = frequency[value]
			mostFrequent = value
		}
	}

	return mostFrequent
}
