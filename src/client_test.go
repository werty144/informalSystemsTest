package main

import (
	"net"
	"testing"
)

func TestPlayGame(t *testing.T) {
	stopChanel := make(chan struct{})
	defer close(stopChanel)

	trueValue := 3
	var addresses []net.TCPAddr
	for i := 0; i < 3; i++ {
		addr := startAgent(trueValue, 10, false, stopChanel)
		addresses = append(addresses, *addr)
	}
	for i := 0; i < 2; i++ {
		addr := startAgent(trueValue, 10, true, stopChanel)
		addresses = append(addresses, *addr)
	}

	err, v := getNetworkValue(addresses)
	if err != nil {
		t.Errorf(err.Error())
	}
	if v != trueValue {
		t.Fail()
	}
}
