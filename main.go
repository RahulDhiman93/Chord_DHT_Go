package main

import (
	"fmt"
	"net/http"
	"sync"
)

var SERVERS map[int]*http.Server

func startNodeServer(wg *sync.WaitGroup, ip string, port int) {
	defer wg.Done()

	node := newNode(ip, port)
	node.join()
	urlShortenerHandler := URLShortener{node}

	mux := http.NewServeMux()
	mux.HandleFunc("/store", urlShortenerHandler.handleStore)
	mux.HandleFunc("/lookup", urlShortenerHandler.handleLookup)
	mux.HandleFunc("/leave", urlShortenerHandler.handleLeave)

	address := fmt.Sprintf("%s:%d", ip, port)
	server := &http.Server{Addr: address, Handler: mux}
	fmt.Printf("Starting server on %s\n", address)
	SERVERS[port] = server
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("\nError starting server on %s: %v\n\n", address, err)
	}
}

func main() {
	var wg sync.WaitGroup
	SERVERS = make(map[int]*http.Server)
	nodes := []int{8000, 8001, 8002, 8003, 8004}
	for _, port := range nodes {
		wg.Add(1)
		go startNodeServer(&wg, "127.0.0.1", port)
	}
	wg.Wait()
}
