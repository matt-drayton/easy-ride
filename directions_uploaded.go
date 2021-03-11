package main

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

// Calculate the distance of A roads within a journey
func calcARoadDistance(s []maps.Route) int {
	dist := 0
	for i, step := range s[0].Legs[0].Steps {
		fmt.Printf("\nStep %d: %d meters, ", i+1, step.Distance.Meters)
		pretty.Println(step.HTMLInstructions) // html instructions

		regexA, _ := regexp.Compile("A([0-9]+)")
		fmt.Println("Driving on A road? ", regexA.MatchString(step.HTMLInstructions))

		// Find steps that have A roads appearing in the instructions
		if regexA.MatchString(step.HTMLInstructions) == true {
			fmt.Println(regexA.FindString(step.HTMLInstructions))
			dist = dist + step.Distance.Meters
		}
	}
	return dist
}

//
func main() {
	// Make sure you insert your API Key to access the Google Directions API
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyAOtrwQjprZOOtFdOMvrSkDSOdBWr3T7PQ"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.DirectionsRequest{
		Region:      "UK",
		Origin:      "University of Exeter",
		Destination: "Crediton",
	}
	route, _, err := c.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	// Print all steps
	pretty.Println(route[0])
	// pretty.Println(route[0].Legs[0].Steps)

	// Distance made on A road
	distA := calcARoadDistance(route)

	// Total distance of the journey in meters
	distTotal := route[0].Legs[0].Distance.Meters
	fmt.Printf("\n*******************************************************************\n")
	fmt.Printf("Total distance from %s to %s: %d m\n", r.Origin, r.Destination, distTotal)
	fmt.Printf("Distance made on A roads: %d m\n", distA) // total distance in meters

	// Check if the majority of the journey is made on A roads
	if distA > distTotal/2 {
		fmt.Println("\nDouble the charge rate!")
	} else {
		fmt.Println("\nKeep the charge rate!")
	}

}
