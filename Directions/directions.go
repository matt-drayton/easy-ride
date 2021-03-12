package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
	"googlemaps.github.io/maps"
)

type Route struct {
	TotalDistance int `json:"totaldistance`
	ARoadDistance int `json:"aroaddistance`
} 

// Calculate the distance of A roads within a journey
func calcARoadDistance(s []maps.Route) int {
	dist := 0
	for _, step := range s[0].Legs[0].Steps {

		regexA, _ := regexp.Compile("A([0-9]+)")

		// Find steps that have A roads appearing in the instructions
		if regexA.MatchString(step.HTMLInstructions) == true {
			dist = dist + step.Distance.Meters
		}
	}
	return dist
}


func getRouteDistanceHelper(origin, destination string) (distance, aRoadDistance int, err error) {
	// Make sure you insert your API Key to access the Google Directions API
	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("MAPS_API_KEY")))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.DirectionsRequest{
		Region:      "UK",
		Origin:      origin,
		Destination: destination,
	}
	route, _, err := c.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
		return 0, 0, err
	}

	// Distance made on A road
	distA := calcARoadDistance(route)

	// Total distance of the journey in meters
	distTotal := route[0].Legs[0].Distance.Meters

	return distTotal, distA, nil
}

func getRouteDistance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	origin := vars["from"]
	destination := vars["to"]

	totalDistance, aRoadDistance, err := getRouteDistanceHelper(origin, destination)

	if err != nil {
		log.Printf("Error: Could not find route between %s and %s : %s", origin, destination, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Could not find route between %s and %s\"}", origin, destination)))
		return
	}

	route := Route{
		TotalDistance: totalDistance,
		ARoadDistance: aRoadDistance,
	}
	log.Printf("Finding distance between %s and %s", origin, destination)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(route)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/directions/{from}/{to}", getRouteDistance).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	log.Println("Starting Directions Service")
	handleRequests()
}