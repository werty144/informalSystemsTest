package main

import (
	"math/rand"
	"net"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestGetNetworkValue(t *testing.T) {
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

func TestGetMode(t *testing.T) {
	values := []int{1, 1, 2, 2, 1, 3}
	mode := getMode(values)
	if mode != 1 {
		t.Fail()
	}
}

func TestGetModeAmbiguous(t *testing.T) {
	values := []int{1, 1, 2, 2, 1, 2, 3}
	mode := getMode(values)
	if mode != 1 && mode != 2 {
		t.Fail()
	}
}

func TestThreadSafeSlice_Append(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var values ThreadSafeSlice
	var wg sync.WaitGroup
	nWorkers := 1000
	wg.Add(nWorkers)

	for i := 0; i < nWorkers; i++ {
		go func(v int) {
			defer wg.Done()
			delay := rand.Intn(10) + 1
			time.Sleep(time.Duration(delay) * time.Millisecond)
			values.Append(v)
		}(i)
	}
	wg.Wait()

	if len(values.slice) != nWorkers {
		t.Fail()
	}

	sort.Ints(values.slice)

	for i := 0; i < nWorkers; i++ {
		if values.slice[i] != i {
			t.Fail()
		}
	}
}
