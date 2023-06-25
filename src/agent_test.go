package main

import (
	"encoding/binary"
	"net"
	"testing"
	"time"
)

func TestCorrectAgent(t *testing.T) {
	stopChanel := make(chan struct{})
	defer close(stopChanel)

	trueValue := 3
	addr := startAgent(trueValue, 10, false, stopChanel)

	conn, err := net.DialTimeout("tcp", addr.String(), 1*time.Second)
	if err != nil {
		t.Fatal("Error:", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("ping"))
	if err != nil {
		t.Fatal("Error sending message:", err)
	}

	response := make([]byte, 1024)
	length, err := conn.Read(response)
	if err != nil {
		t.Errorf("Error receiving response from "+addr.String()+".", err)
		return
	}
	v := int(binary.BigEndian.Uint32(response[:length]))

	if v != trueValue {
		t.Fail()
	}
}

func TestLiarAgent(t *testing.T) {
	stopChanel := make(chan struct{})
	defer close(stopChanel)

	trueValue := 3
	maxV := 10
	addr := startAgent(trueValue, maxV, true, stopChanel)

	conn, err := net.DialTimeout("tcp", addr.String(), 1*time.Second)
	if err != nil {
		t.Fatal("Error:", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("ping"))
	if err != nil {
		t.Fatal("Error sending message:", err)
	}

	response := make([]byte, 1024)
	length, err := conn.Read(response)
	if err != nil {
		t.Errorf("Error receiving response from "+addr.String()+".", err)
		return
	}
	v := int(binary.BigEndian.Uint32(response[:length]))

	if v == trueValue || v < 1 || v > maxV {
		t.Fail()
	}
}
