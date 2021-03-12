package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type cheapestDriverResponse struct {
	CheapestDriver driver `json:"cheapest_driver"`
	NoOfDrivers int `json:"no_of_drivers"`
}

type driver struct {
	Username string `json:"username"`
	Name string `json:"name"`
	Rate int `json:"rate"`
}

type route struct {
	TotalDistance int `json:"totaldistance`
	ARoadDistance int `json:"aroaddistance`
} 

type journey struct {
	StartPoint string `json:"start_point"`
	EndPoint string `json:"end_point"`
	TotalDistance int `json:"total_distance"`
	ARoadDistance int `json:"a_road_distance"`
	BestDriver driver `json:"best_driver"`
	Cost int `json:"cost"`
}

func getJourney(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	origin := vars["from"]
	destination := vars["to"]

	// Get route distance
	resp, err := http.Get(fmt.Sprintf("http://directions-service:8000/directions/%s/%s", origin, destination))

	if err != nil {

	}

	var distances route
	json.NewDecoder(resp.Body).Decode(&distances)

	// Get cheapest driver
	resp, err = http.Get("http://roster-service:8000/roster")
	var fetchedCheapestDriverResponse cheapestDriverResponse
	json.NewDecoder(r.Body).Decode(&fetchedCheapestDriverResponse)


	cost := calculateCost(distances, fetchedCheapestDriverResponse)

	response := journey {
		StartPoint: origin,
		EndPoint: destination,
		TotalDistance: distances.TotalDistance,
		ARoadDistance: distances.ARoadDistance,
		BestDriver: fetchedCheapestDriverResponse.CheapestDriver,
		Cost: cost,
	}

	log.Println(fmt.Sprintf("Journey between %s and %s calculated at %dp with driver %s", origin, destination, cost, 
																						  response.BestDriver.Username))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func calculateCost(routeDetails route, driverDetails cheapestDriverResponse) int {
	cheapestDriverDetails := driverDetails.CheapestDriver
	noOfDrivers := driverDetails.NoOfDrivers

	cost := cheapestDriverDetails.Rate * routeDetails.TotalDistance

	// Check if over half of distance is on A-Roads.
	// Integer division is ideal here. No need to convert types.
	if routeDetails.ARoadDistance > (routeDetails.TotalDistance / 2) {
		cost *= 2
	}

	if noOfDrivers < 5 {
		cost *= 2
	}

	return cost
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/journey/{from}/{to}", getJourney).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	log.Println("Starting Journey Service")
	handleRequests()
}