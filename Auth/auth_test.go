package main

import (
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

type driver struct {
	Username string `json:"username"`
	Name string `json:"name"`
	Rate int `json:"rate"`
}

func TestAuth(t *testing.T) {
	// Sleep to ensure auth service has finished build
	time.Sleep(3 * time.Second)
	// Valid Credentials
	data := url.Values{}
	data.Set("username", "sebvet")
	data.Set("password", "astonmartin")
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://auth-service:8000/login", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(req) 

	if err != nil || resp.StatusCode != 200 {
		log.Println("Failed valid credentials sign-in unit test.")
		log.Println(err)
		log.Println(resp.StatusCode)
		t.Fail()
	}
	var token Token
	json.NewDecoder(resp.Body).Decode(&token)
	
	// Invalid Credentials
	data = url.Values{}
	data.Set("username", "fakename")
	data.Set("password", "fakepassword")
	req, err = http.NewRequest("POST", "http://auth-service:8000/login", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err = client.Do(req) 
	if resp.StatusCode == 200 {
		log.Println(err)
		log.Println(resp.StatusCode)
		log.Println("Failed invalid credentials sign-in unit test.")
		t.Fail()
	}

	// Take token from valid login and verify
	r, err := http.Get("http://auth-service:8000/validate/"+token.Token)
	
	if err != nil || r.StatusCode != http.StatusOK {
		log.Println("Failed to validate correct JWT")
		t.Fail()
	}

	var authenticatedDriver driver
	expectedDriver := driver {
		Username: "sebvet",
		Name: "Sebastian Vettel",
		Rate: 0,
	}
	json.NewDecoder(r.Body).Decode(&authenticatedDriver)
	
	if authenticatedDriver != expectedDriver {
		log.Println("Failed to fetch correct user data from JWT")
		t.Fail()
	}

	r, err = http.Get("http://auth-service:8000/validate/"+"RANDOMNOTVALIDTOKEN")
	
	if r.StatusCode == http.StatusOK {
		log.Println("Failed to correctly reject invalid token")
		t.Fail()
	}

}
