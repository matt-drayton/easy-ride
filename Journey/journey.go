package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

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
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	origin := vars["from"]
	destination := vars["to"]

	// Get route distance
	resp, err := http.Get(fmt.Sprintf("http://directions-service:8000/directions/%s/%s", origin, destination))

	if err != nil {
		log.Printf("Error: Could not fetch route between %s and %s : %s", origin, destination, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Could not fetch route between %s and %s.\"}", origin, destination)))
		return
	}

	var distances route
	json.NewDecoder(resp.Body).Decode(&distances)

	// Get cheapest driver
	resp, err = http.Get("http://roster-service:8000/roster")
	if err != nil {
		log.Printf("Error fecthing roster: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\": \"Could not fetch roster data\"}"))
		return
	}
	
	var fetchedDrivers []driver 
	json.NewDecoder(resp.Body).Decode(&fetchedDrivers)

	cheapestDriver := getCheapestDriver(fetchedDrivers)

	cost := calculateCost(distances, fetchedDrivers)

	response := journey {
		StartPoint: origin,
		EndPoint: destination,
		TotalDistance: distances.TotalDistance,
		ARoadDistance: distances.ARoadDistance,
		BestDriver: cheapestDriver,
		Cost: cost,
	}

	log.Println(fmt.Sprintf("Journey between %s and %s calculated at %dp with driver %s", origin, destination, cost, 
																						  response.BestDriver.Username))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func calculateCost(routeDetails route, availableDrivers []driver) int {
	cheapestDriver := getCheapestDriver(availableDrivers)
	noOfDrivers := len(availableDrivers)

	cost := cheapestDriver.Rate * (routeDetails.TotalDistance / 1000)

	// Check if over half of distance is on A-Roads.
	// Integer division is ideal here. No need to convert types.
	if routeDetails.ARoadDistance > (routeDetails.TotalDistance / 2) {
		cost *= 2
	}

	if noOfDrivers < 5 {
		cost *= 2
	}

	currentHour, _, _ := time.Now().Clock()

	if currentHour >= 23 || currentHour <= 6 {
		cost *= 2
	}

	return cost
}

func getCheapestDriver(drivers []driver) driver{

	lowestRate := -1
	var lowestDriver driver

	// Find the driver with the lowest rate. This will always be the best driver for the route.
	for _, currentDriver := range drivers {
		// If first pass, initialise. Lowest rate can never naturally be -1
		if lowestRate == -1 {
			lowestRate = currentDriver.Rate
			lowestDriver = currentDriver
			continue
		}

		if currentDriver.Rate < lowestRate {
			lowestDriver = currentDriver
			lowestRate = currentDriver.Rate
		}
	}
	return lowestDriver
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