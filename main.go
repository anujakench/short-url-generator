package main

import (
	base64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ShortURL struct {
	urls map[string]string
}

type ReqBody struct {
	URL string
}

func main() {
	for {
		shorturl := &ShortURL{
			urls: make(map[string]string),
		}

		http.HandleFunc("/shortenurl", shorturl.HandleShortenURL)
		http.HandleFunc("/deleteurl", shorturl.HandleDeleteURL)
		http.HandleFunc("/redirecturl", shorturl.HandleRedirectURL)

		// Start the HTTP server. Listen for incoming requests
		fmt.Println("URL shortening service has started.")
		http.ListenAndServe(":8080", nil)
	}
}

func (s *ShortURL) HandleShortenURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered HandleShortenURL")
	fmt.Printf("%s\n", s.urls)
	// check if HTTP request type is POST
	var longURL ReqBody
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// check if HTTP request has a body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body: ", http.StatusBadRequest)
		fmt.Printf("server: could not read request body: %s\n", err)
	}
	fmt.Printf("server: request body: %s\n", reqBody)
	// Unmarshall the data
	err = json.Unmarshal(reqBody, &longURL)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusInternalServerError)
		fmt.Println("Error unmarshalling JSON: ", err)
	}
	// Fetch the long URL
	if longURL.URL == "" {
		http.Error(w, "URL data is missing", http.StatusBadRequest)
		return
	}
	fmt.Printf("longURL %s\n", longURL.URL)
	// Create short URL
	encodedStr := encodeOriginalURL(longURL.URL)
	shortURL := fmt.Sprintf("http://myshorturl.com/%s", encodedStr)
	s.urls[shortURL] = longURL.URL
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error happened in JSON marshal. Err: %s", http.StatusInternalServerError)
	}
	w.Write(jsonResp)
}

func (s *ShortURL) HandleDeleteURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered HandleDeleteURL")
	fmt.Printf("%s\n", s.urls)
	// check if HTTP request type is PUT
	var shortURL ReqBody
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// check if HTTP request has a body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body: ", http.StatusBadRequest)
		fmt.Printf("server: could not read request body: %s\n", err)
	}
	fmt.Printf("server: request body: %s\n", reqBody)
	// Unmarshall the data
	err = json.Unmarshal(reqBody, &shortURL)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusInternalServerError)
		fmt.Println("Error unmarshalling JSON: ", err)
	}
	// Fetch the long URL
	if shortURL.URL == "" {
		http.Error(w, "URL data is missing", http.StatusBadRequest)
		return
	}
	fmt.Printf("shortURL %s\n", shortURL.URL)
	_, ok := s.urls[shortURL.URL]
	if ok {
		delete(s.urls, shortURL.URL)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Status OK"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Error happened in JSON marshal. Err: %s", http.StatusInternalServerError)
		}
		w.Write(jsonResp)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Status Bad Request"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Error happened in JSON marshal. Err: %s", http.StatusBadRequest)
		}
		w.Write(jsonResp)
		return
	}
}

func (s *ShortURL) HandleRedirectURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered HandleRedirectURL")
	fmt.Printf("%s\n", s.urls)
	// check if HTTP request type is POST
	var shortURL ReqBody
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// check if HTTP request has a body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body: ", http.StatusBadRequest)
		fmt.Printf("server: could not read request body: %s\n", err)
	}
	fmt.Printf("server: request body: %s\n", reqBody)
	// Unmarshall the data
	err = json.Unmarshal(reqBody, &shortURL)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusInternalServerError)
		fmt.Println("Error unmarshalling JSON: ", err)
	}
	// Fetch the long URL
	if shortURL.URL == "" {
		http.Error(w, "URL data is missing", http.StatusBadRequest)
		return
	}
	fmt.Printf("shortURL %s\n", shortURL.URL)
	longurl, found := s.urls[shortURL.URL]
	if !found {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Redirect user to long URL
	http.Redirect(w, r, longurl, http.StatusMovedPermanently)
}

// Private function to generate a base64 encoded string
func encodeOriginalURL(url string) string {
	encodedString := base64.StdEncoding.EncodeToString([]byte(url))
	fmt.Println(encodedString)
	return encodedString
}
