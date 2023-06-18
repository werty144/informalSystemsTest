package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func startServer(stopChanel <-chan struct{}) {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
		}
	}(listener)

	go func() {
		for {
			select {
			case <-stopChanel:
				err := listener.Close()
				if err != nil {
					//log.Fatal("Error closing listener:", err)
				}
				return
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			//log.Println("Error accepting connection:", err)
			return
		}

		go handleConnection(conn, stopChanel) // Handle each connection in a separate goroutine
	}
}

func handleConnection(conn net.Conn, stopChanel <-chan struct{}) {
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
		_, err = conn.Write(buffer[:n])
		if err != nil {
			//log.Println("Error writing to connection:", err)
			return
		}
	}
}
