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
	/*
		Tests getting a network value.
	*/

	// create a stop chanel for agents
	stopChanel := make(chan struct{})
	defer close(stopChanel)

	trueValue := 3
	var addresses []net.TCPAddr

	// create correct agents
	for i := 0; i < 3; i++ {
		addr := startAgent(trueValue, 10, false, stopChanel)
		addresses = append(addresses, *addr)
	}

	//create liars
	for i := 0; i < 2; i++ {
		addr := startAgent(trueValue, 10, true, stopChanel)
		addresses = append(addresses, *addr)
	}

	// use the getNetworkValue function to get the network value
	err, v := getNetworkValue(addresses)
	if err != nil {
		t.Errorf(err.Error())
	}

	// check whether the returned value is indeed a network value
	if v != trueValue {
		t.Fail()
	}
}

func TestGetMode(t *testing.T) {
	/*
		Tests the get mode function in case there is a single mode in the collection.
	*/

	values := []int{1, 1, 2, 2, 1, 3} // create the test collection
	mode := getMode(values)           // call the getMode function

	// Check if the returned value is indeed the most frequent one
	if mode != 1 {
		t.Fail()
	}
}

func TestGetModeAmbiguous(t *testing.T) {
	/*
		Tests the getMode function in case there are several modes
	*/

	values := []int{1, 1, 2, 2, 1, 2, 3} // create a collection with several modes
	mode := getMode(values)              // call getMode function to get the mode

	// Check if the returned result is indeed one of the modes
	if mode != 1 && mode != 2 {
		t.Fail()
	}
}

func TestThreadSafeSlice_Append(t *testing.T) {
	/*
		Tests the ThreadSafeSlice structure by concurrently appending elements to it.
	*/

	rand.Seed(time.Now().UnixNano()) // make sure randomness is fresh each invocation
	var values ThreadSafeSlice       // create the structure
	var wg sync.WaitGroup            // create the synchronization group
	nWorkers := 1000
	wg.Add(nWorkers)

	// launch concurrent goroutines adding values to the structure with small random delays
	for i := 0; i < nWorkers; i++ {
		go func(v int) {
			defer wg.Done()
			delay := rand.Intn(10) + 1
			time.Sleep(time.Duration(delay) * time.Millisecond)
			values.Append(v)
		}(i)
	}
	wg.Wait()

	// Check if the number of elements in the structure corresponds to the number of performed additions
	if len(values.slice) != nWorkers {
		t.Fail()
	}

	// Check if all elements that were added are actually present in the structure
	sort.Ints(values.slice)
	for i := 0; i < nWorkers; i++ {
		if values.slice[i] != i {
			t.Fail()
		}
	}
}
