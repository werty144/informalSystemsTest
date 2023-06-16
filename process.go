package main

import (
	"fmt"
	"time"
)

func echo(stopChanel <-chan struct{}) {
	for {
		select {
		default:
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Guy")
		case <-stopChanel:
			fmt.Println("Stopped")
			return
		}
	}
}
