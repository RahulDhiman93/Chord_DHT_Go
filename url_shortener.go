package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

const shortKeyLength = 8

type URLShortener struct {
	node *Node
}

func NewURLShortener(node *Node) *URLShortener {
	return &URLShortener{node: node}
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
