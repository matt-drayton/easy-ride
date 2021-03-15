package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

type Token struct {
	Token string `json:"token"`
}

type rosterReq struct {
	Token string `json:"token"`
	Rate int `json:"rate"`
}

func TestRoster(t *testing.T) {
	// Sleep to ensure other services have finished build
	time.Sleep(3 * time.Second)

	// Valid Credentials. We can assume this works. If it doesn't, it'll be caught in auth_test.go
	data := url.Values{}
	data.Set("username", "sebvet")
	data.Set("password", "astonmartin")
	client := &http.Client{}

	req, _ := http.NewRequest("POST", "http://auth-service:8000/login", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(req) 

	var token Token
	json.NewDecoder(resp.Body).Decode(&token)

	// Test getting roster. Should be empty. 
	resp, err := http.Get("http://roster-service:8000/roster")
	if err != nil {
		log.Println("Failed fetching roster")
		t.Fail()
	}
	
	var fetchedDrivers []driver 
	json.NewDecoder(resp.Body).Decode(&fetchedDrivers)

	if len(fetchedDrivers) != 0 {
		log.Println("Failed fetching roster when it should be empty")
		t.Fail()
	}

	// Try joining roster
	joinRosterReq := rosterReq{
		Token: token.Token,
		Rate: 5,
	}
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(joinRosterReq)
	resp, err = http.Post("http://roster-service:8000/roster", "application/json", payload)

	var returnedDriver driver
	json.NewDecoder(resp.Body).Decode(&returnedDriver)

	if returnedDriver.Rate != 5 {
		log.Println("Failed to set rate correctly")
		t.Fail()
	}

	// Test updating rate
	joinRosterReq.Rate = 10
	payload = new(bytes.Buffer)
	json.NewEncoder(payload).Encode(joinRosterReq)
	req, _ = http.NewRequest("PUT", "http://roster-service:8000/roster", payload)
	req.Header.Add("Content-Type", "application/json")
	resp, _ = client.Do(req) 

	json.NewDecoder(resp.Body).Decode(&returnedDriver)

	if returnedDriver.Rate != 10 {
		log.Println("Failed to update rate correctly")
		t.Fail()
	}

	// Test leaving roster
	payload = new(bytes.Buffer)
	json.NewEncoder(payload).Encode(token)
	req, _ = http.NewRequest("DELETE", "http://roster-service:8000/roster", payload)
	req.Header.Add("Content-Type", "application/json")
	resp, _ = client.Do(req) 

	// Get users in roster to see if removal has worked
	resp, err = http.Get("http://roster-service:8000/roster")
	if err != nil {
		log.Println("Failed fetching roster")
		t.Fail()
	}
	
	json.NewDecoder(resp.Body).Decode(&fetchedDrivers)

	if resp.StatusCode != 200 || len(fetchedDrivers) != 0{
		log.Println("Failed to leave roster")
		t.Fail()
	}
}	
