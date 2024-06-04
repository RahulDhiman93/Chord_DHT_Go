package main

import (
	"fmt"
	"net/http"
	"sync"
)

func startNodeServer(wg *sync.WaitGroup, ip string, port int) {
	defer wg.Done()

	node := newNode(ip, port)
	node.join()
	urlShortenerHandler := URLShortener{node}

	mux := http.NewServeMux()
	mux.HandleFunc("/store", urlShortenerHandler.handleStore)
	mux.HandleFunc("/lookup", urlShortenerHandler.handleLookup)

	address := fmt.Sprintf("%s:%d", ip, port)
	fmt.Printf("Starting server on %s\n", address)
	err := http.ListenAndServe(address, mux)
	if err != nil {
		fmt.Printf("Error starting server on %s: %v\n", address, err)
	}
}

func main() {
	var wg sync.WaitGroup

	nodes := []int{8000, 8001, 8002, 8003, 8004}
	for _, port := range nodes {
		wg.Add(1)
		go startNodeServer(&wg, "127.0.0.1", port)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Assume newNode returns a *Node
		node1 := newNode("127.0.0.1", 8000)
		node2 := newNode("127.0.0.1", 8001)
		node3 := newNode("127.0.0.1", 8002)
		node4 := newNode("127.0.0.1", 8003)
		node5 := newNode("127.0.0.1", 8004)

		node1.join()
		node2.join()
		node3.join()
		node4.join()
		node5.join()

		for i := 0; i < 3; i++ {
			node1.stabilize()
			node2.stabilize()
			node3.stabilize()
			node4.stabilize()
			node5.stabilize()
			node1.fixFingers()
			node2.fixFingers()
			node3.fixFingers()
			node4.fixFingers()
			node5.fixFingers()
		}

		node1.printNodeData()
		node2.printNodeData()
		node3.printNodeData()
		node4.printNodeData()
		node5.printNodeData()
	}()

	// Wait for all goroutines to finish
	wg.Wait()
}
