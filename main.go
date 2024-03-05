package main

import (
	base64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ShortURL struct {
	url   string
	hours [168]uint64 // URL access time on hourly basis
}

type ShortURLS struct {
	urls map[ShortURL]string
}

type ReqBody struct {
	URL string
}

type AccessedReqBody struct {
	Url        string
	AccessTime string
}

var EpochTime time.Time

func main() {
	// Calculate the epoch time
	EpochTime = time.Now()
	for {
		shorturl := &ShortURLS{
			urls: make(map[ShortURL]string),
		}

		http.HandleFunc("/shortenurl", shorturl.HandleShortenURL)
		http.HandleFunc("/deleteurl", shorturl.HandleDeleteURL)
		http.HandleFunc("/redirecturl", shorturl.HandleRedirectURL)
		http.HandleFunc("/urlaccessed", shorturl.HandleAccessedTimeURL)

		// Start the HTTP server. Listen for incoming requests
		fmt.Println("URL shortening service has started.")
		http.ListenAndServe(":8080", nil)
	}
}

func (s *ShortURLS) HandleShortenURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered HandleShortenURL")
	currentTime := time.Now()
	difference := currentTime.Sub(EpochTime)
	fmt.Printf("hours %f\n", difference.Hours())
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
	fmt.Printf("longURL %s\n", longURL.URL)
	// Fetch the long URL
	if longURL.URL == "" {
		http.Error(w, "URL data is missing", http.StatusBadRequest)
		return
	}
	fmt.Printf("longURL %s\n", longURL.URL)
	// Create short URL
	var str ShortURL
	encodedStr := encodeOriginalURL(longURL.URL)
	shortUrl := fmt.Sprintf("http://myshorturl.com/%s", encodedStr)
	str.url = shortUrl
	if existsShortUrl(str.url, s) {
		http.Error(w, "Duplicate short URL found,", http.StatusConflict)
		return
	}
	fmt.Printf("hours %d\n", int64(difference.Hours()))
	str.hours[int64(difference.Hours())] += 1
	// Store the short URL in hash map

	s.urls[str] = longURL.URL
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error happened in JSON marshal. Err: %s", http.StatusInternalServerError)
	}
	w.Write(jsonResp)
	fmt.Println("URL's")
	fmt.Printf("%+v\n", s.urls)
}

func (s *ShortURLS) HandleDeleteURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered HandleDeleteURL")
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
	for shortUrl, _ := range s.urls {
		if shortUrl.url == shortURL.URL {
			delete(s.urls, shortUrl)
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
	fmt.Println("URL's")
	fmt.Printf("%+v\n", s.urls)
}

func (s *ShortURLS) HandleRedirectURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered HandleRedirectURL")
	currentTime := time.Now()
	difference := currentTime.Sub(EpochTime)
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
	for short, long := range s.urls {
		fmt.Printf("short.url: %s\t", short.url)
		fmt.Printf("shortURL.URL: %s\n", shortURL.URL)
		if short.url == shortURL.URL {
			fmt.Printf("hours before %d\n", short.hours[int64(difference.Hours())])
			short_copy := short
			delete(s.urls, short)
			short_copy.hours[int64(difference.Hours())] += 1
			fmt.Printf("hours after %d\n", short_copy.hours[int64(difference.Hours())])
			s.urls[short_copy] = long
			// Redirect user to long URL
			w.Header().Set("Content-Type", "application/json")
			http.Redirect(w, r, long, http.StatusMovedPermanently)
			return
		}
	}
	fmt.Println("URL's")
	fmt.Printf("%+v\n", s.urls)
	http.Error(w, "Short URL not found", http.StatusNotFound)
}

func (s *ShortURLS) HandleAccessedTimeURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered HandleAccessedTimeURL")
	// check if HTTP request type is GET
	var shortURL AccessedReqBody
	if r.Method != http.MethodGet {
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
	fmt.Printf("shortURL %+v\n", shortURL)
	accessTime := shortURL.AccessTime
	var hours [168]uint64
	for short, _ := range s.urls {
		if short.url == shortURL.Url {
			hours = short.hours
			break
		}
	}
	format := strings.Split(accessTime, " ")
	// validate the current time and access request time
	currentTime := time.Now()
	var hrs int
	hrs = int(currentTime.Sub(EpochTime).Hours())
	intHrs, err := strconv.Atoi(format[0])
	if err != nil {
		fmt.Printf("Error %+v\n", err)
		return
	}
	fmt.Printf("hrs %d, intHrs %d\n", hrs, intHrs)
	if currentTime.Before(EpochTime) || hrs < intHrs {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	fmt.Printf("format %+v\n", format)
	var urlTime uint64
	if format[1] == "hours" {
		hrs, err = strconv.Atoi(format[0])
		if err != nil {
			fmt.Printf("Error converting URL access time\n")
			return
		}
		for i := 0; i <= hrs; i++ {
			urlTime += hours[i]
		}
	} else if format[1] == "weeks" {
		weeks, err := strconv.Atoi(format[0])
		var week_hrs int
		if err != nil {
			fmt.Printf("Error converting URL access time\n")
			return
		}
		if weeks != 0 {
			week_hrs = weeks * 168
		}
		for i := 0; i <= week_hrs; i++ {
			urlTime += hours[i]
		}
	}
	fmt.Printf("url time %d\n", urlTime)
	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
        <h2>URL Shortener</h2>
        <p>Shortened URL: <%s></p>
		<p>Access Time: <a href="%s is accessed %d times since %s %s.</a></p>
        <form method="get" action="/accessedURL">
            
        </form>
    `, shortURL.Url, shortURL.Url, urlTime, format[0], format[1])
	fmt.Fprintln(w, responseHTML)
}

// Private function to generate a base64 encoded string
func encodeOriginalURL(url string) string {
	encodedString := base64.StdEncoding.EncodeToString([]byte(url))
	fmt.Println(encodedString)
	return encodedString
}

// Function to check if duplicate short URL exists
func existsShortUrl(url string, s *ShortURLS) bool {
	for short, _ := range s.urls {
		if short.url == url {
			return true
		}
	}
	return false
}
