package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

func startServer(v int, maxV int, isLiar bool, stopChanel <-chan struct{}) {
	// Start listening to TCP
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}

	//Defer closing listener with error handling
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
		}
	}(listener)

	port := listener.Addr().(*net.TCPAddr).Port
	log.Println(port, isLiar)

	// Periodically check if it's needed to quite the server
	go func() {
		for {
			select {
			case <-stopChanel:
				err := listener.Close()
				if err != nil {
				}
				return
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	processValue := getProcessValue(v, maxV, isLiar)

	// Accept and handle new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		go handleConnection(processValue, conn, stopChanel) // Handle each connection in a separate goroutine
	}
}

func handleConnection(processValue int, conn net.Conn, stopChanel <-chan struct{}) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			//log.Fatal("Error closing the connection", err)
		}
	}(conn)

	go func() {
		for {
			select {
			case <-stopChanel:
				err := conn.Close()
				if err != nil {
				}
				return
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	fmt.Println("New connection established:", conn.RemoteAddr())

	// Example: Echo server - read data from the client and send it back
	buffer := make([]byte, 1024)
	data := make([]byte, 4) // Assuming 32-bit integer
	binary.BigEndian.PutUint32(data, uint32(processValue))
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			//log.Println("Error reading from connection:", err)
			return
		}
		if n == 0 {
			return
		}

		// Echo back the received data
		_, err = conn.Write(data)
		if err != nil {
			//log.Println("Error writing to connection:", err)
			return
		}
	}
}

func getProcessValue(v int, maxV int, isLiar bool) int {
	// If process is honest, return v
	if !isLiar {
		return v
	}

	// Sample fake value uniformly from [1, maxV] \ {v}
	fakeValue := rand.Intn(maxV-1) + 1
	if fakeValue >= v {
		fakeValue++
	}

	return fakeValue
}