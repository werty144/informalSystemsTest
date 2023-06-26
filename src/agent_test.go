package main

import (
	"encoding/binary"
	"net"
	"testing"
	"time"
)

func TestCorrectAgent(t *testing.T) {
	/*
		Tests the creation and the response of the correct agent
	*/

	// create a stop chanel for the agent
	stopChanel := make(chan struct{})
	defer close(stopChanel)

	// start the agent
	trueValue := 3
	addr := startAgent(trueValue, 10, false, stopChanel)

	// connect to the agent
	conn, err := net.DialTimeout("tcp", addr.String(), 1*time.Second)
	if err != nil {
		t.Fatal("Error:", err)
	}
	defer conn.Close()

	// query the agent
	_, err = conn.Write([]byte("ping"))
	if err != nil {
		t.Fatal("Error sending message:", err)
	}

	// get agents response
	response := make([]byte, 1024)
	length, err := conn.Read(response)
	if err != nil {
		t.Errorf("Error receiving response from "+addr.String()+".", err)
		return
	}
	v := int(binary.BigEndian.Uint32(response[:length]))

	// check if the received value is a value the agent was intended to have
	if v != trueValue {
		t.Fail()
	}
}

func TestLiarAgent(t *testing.T) {
	/*
		Tests the creation and the response of the liar agent
	*/

	// create a stop chanel for the agent
	stopChanel := make(chan struct{})
	defer close(stopChanel)

	// create the agent
	trueValue := 3
	maxV := 10
	addr := startAgent(trueValue, maxV, true, stopChanel)

	// connect to the agent
	conn, err := net.DialTimeout("tcp", addr.String(), 1*time.Second)
	if err != nil {
		t.Fatal("Error:", err)
	}
	defer conn.Close()

	// query the agent
	_, err = conn.Write([]byte("ping"))
	if err != nil {
		t.Fatal("Error sending message:", err)
	}

	// get the agents response
	response := make([]byte, 1024)
	length, err := conn.Read(response)
	if err != nil {
		t.Errorf("Error receiving response from "+addr.String()+".", err)
		return
	}
	v := int(binary.BigEndian.Uint32(response[:length]))

	// Check that the returned value is not the network value and that it is within the allowed range
	if v == trueValue || v < 1 || v > maxV {
		t.Fail()
	}
}
