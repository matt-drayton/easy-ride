package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type driver struct {
	Username string `json:"username"`
	Name string `json:"name"`
	Rate int `json:"rate"`
}


type driverRateRequest struct {
	Token string `json:"token"`
	Rate int `json:"rate"`
}

var Roster = map[string]driver{}


func authenticateUser(token string) (*driver, error) {

	r, err := http.Get("http://auth-service:8000/validate/"+token)
	
	if err != nil || r.StatusCode != http.StatusOK {
		return nil, errors.New("unauthorised jwt")
	}

	var authenticatedDriver driver
	json.NewDecoder(r.Body).Decode(&authenticatedDriver)

	// Note that just because driver is authenticated, doesn't mean they are in roster
	// Catch on other side
	return &authenticatedDriver, nil
}

// Requires authentication
func joinRoster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	
	if err != nil {
		log.Println("Error: Parsing request to join roster failed.")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\": \"Parsing request to join roster failed\"}"))
		return
	}

	var requestData driverRateRequest
	err = json.Unmarshal(body, &requestData)
	
	if err != nil {
		log.Println("Error: Request is missing JWT token or rate.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"Request is missing JWT token or rate\"}"))
		return
	}

	user, err := authenticateUser(requestData.Token)

	if err != nil {
		log.Println("Error: Invalid JWT token.")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{\"error\": \"Invalid JWT token\"}"))
		return
	}

	_, ok := Roster[user.Username]

	// Check if driver is already in roster.
	if ok {
		log.Println("Error: User is already in roster.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"User is already in roster\"}"))
		return
	}

	// Cannot have a rate of less than or equal to 0p.
	if requestData.Rate <= 0 {
		log.Println("Error: Invalid rate value supplied.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"Invalid rate value supplied\"}"))
		return
	}

	// At this point, we can safely add driver to the roster with the given rate
	// Note we do not ask for username or name in this endpoint.
	// By the time they have a token, they have already given this information.

	user.Rate = requestData.Rate
	Roster[user.Username] = *user

	log.Printf("User %s added to roster with rate %dp", user.Username, user.Rate)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Requires authentication
func leaveRoster(w http.ResponseWriter, r *http.Request) {

}

// Requires authentication
func changeRate(w http.ResponseWriter, r *http.Request) {

}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/roster", joinRoster).Methods("POST")
	router.HandleFunc("/roster", leaveRoster).Methods("DELETE")
	router.HandleFunc("/roster", changeRate).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	log.Println("Starting Roster Service")
	handleRequests()
}