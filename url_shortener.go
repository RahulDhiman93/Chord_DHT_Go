package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const shortKeyLength = 8

type URLShortener struct {
	node *Node
}

func (us *URLShortener) generateShortKey() string {
	rand.Seed(time.Now().UnixNano())
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, shortKeyLength)
	for i := 0; i < shortKeyLength; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}
	return string(result)
}

func (us *URLShortener) shortenURL(longURL string) string {
	shortKey := us.generateShortKey()
	us.node.store(shortKey, longURL)
	return shortKey
}

func (us *URLShortener) getLongURL(shortKey string) (string, bool) {
	return us.node.lookup(shortKey)
}

func (us *URLShortener) handleStore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	longURL := requestBody.URL
	if longURL == "" {
		http.Error(w, "url parameter is required", http.StatusBadRequest)
		return
	}

	shortKey := us.shortenURL(longURL)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"shortKey": shortKey})
}

func (us *URLShortener) handleLookup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	shortKey := requestBody.Key
	if shortKey == "" {
		http.Error(w, "key parameter is required", http.StatusBadRequest)
		return
	}

	longURL, ok := us.getLongURL(shortKey)
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"longURL": longURL})
}

func (us *URLShortener) handleLeave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Port int `json:"port"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	port := requestBody.Port

	server, exists := SERVERS[port]
	if !exists {
		http.Error(w, "Server not found on the specified port", http.StatusNotFound)
		return
	}

	node := us.node
	if node.port != port {
		http.Error(w, "Node not associated with the specified port", http.StatusBadRequest)
		return
	}

	node.leave()
	fmt.Printf("Node %d on port %d is leaving the network", node.id, port)

	go func() {
		time.Sleep(1 * time.Second)
		err := server.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("Node %d got error %d while shutting down", node.id, err)
			return
		}

		fmt.Println("STABLING AND FIXING FINGER TABLE ENTRIES")
		fmt.Println("GLOBAL NODES STATE:")
		for _, v := range NETWORK_NODES {
			fmt.Println(v.port)
		}
		fmt.Println("\nUPDATED NODE DATA----->")
		for _, v := range NETWORK_NODES {
			v.printNodeData()
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}
