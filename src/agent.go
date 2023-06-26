package main

import (
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"time"
)

func startAgent(networkValue int, maxV int, isLiar bool, stopChanel <-chan struct{}) *net.TCPAddr {
	/*
		Starts the agent process that listens for the incoming TCP connections.
		The agent sends its value as a response to any received message.
		Arguments:
			networkValue: network value
			maxV: maximal value of the process
			isLiar: boolean telling whether the agent is a liar
			stopChanel: pointer to the chanel used to stop the agent
		Returns:
			TCP address of created agent
	*/

	// Start listening to TCP picking the random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal("Error starting the agent server:", err)
	}

	// Periodically check if it's needed to quit the agent
	go func() {
		for {
			select {
			case <-stopChanel:
				// Closing the listener would result in the interruption of listener.Accept
				// and would lead to the graceful termination of the agent
				// freeing acquired resources
				err := listener.Close()
				if err != nil {
				}
				return
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	// get the value of an agent based on the network value and the isLiar status
	agentValue := getAgentValue(networkValue, maxV, isLiar)

	// Keep accepting connections until listener is not closed
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			go handleConnection(agentValue, conn, stopChanel) // Handle each connection in a separate goroutine
		}
	}()

	return listener.Addr().(*net.TCPAddr)
}

func handleConnection(agentValue int, conn net.Conn, stopChanel <-chan struct{}) {
	/*
		Handles individual connections to the agent.
		Replies with the agent value to any received message.
		Arguments:
			agentValue: value of the agent
			conn: connection to handle
			stopChanel: chanel to be informed about the need to close the connection
	*/
	defer conn.Close()

	// Periodically check if the connection needs to be closed
	go func() {
		for {
			select {
			case <-stopChanel:
				// Closing the connection would result in the interruption of conn.Read
				// and would lead to the graceful termination of the agent
				// freeing acquired resources
				err := conn.Close()
				if err != nil {
				}
				return
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	buffer := make([]byte, 1024) // assuming the incoming message size is <= 1024 bytes
	data := make([]byte, 4)      // Assuming the process value is 32-bit integer
	binary.BigEndian.PutUint32(data, uint32(agentValue))
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		if n == 0 {
			return
		}

		_, err = conn.Write(data) // reply with the agent value
		if err != nil {
			return
		}
	}
}

func getAgentValue(networkVale int, maxV int, isLiar bool) int {
	/*
		Computes the agent value.
		If agent is not a liar, simpy returns network value.
		If it is returns x: 1 <= x <= maxV, x != network value
		Arguments:
			networkValue: network value
			maxV: maximal possible agent value
			isLiar: boolean telling whether the agent is a liar
		Returns:
			Agent value
		If maxV <= 1, no guarantees are given
	*/

	// If the process is correct, return networkVale
	if !isLiar {
		return networkVale
	}

	// Sample fake value uniformly from [1, maxV] \ {networkVale}
	fakeValue := rand.Intn(maxV-1) + 1
	if fakeValue >= networkVale {
		fakeValue++
	}

	return fakeValue
}
