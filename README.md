# Easy-Ride

EasyRide is a fictional UK based ride-sharing app powered by microservices written in Go and deployed using Docker.

## API Documentation

Full API documentation can be found in the `API Specification` folder. The API is described by `.yml` folders that have been exported to `html` using Redoc. 

The API documentation includes an internal and external address for each service. The internal address is used for microservices to communicate with each other. The external address is used for requests from outside of the application. In this instance, external addresses are all on http://localhost with varying ports.

## Services

The EasyRide service has been broken into the following microservices, each of which resides in its own directory:

- Auth
  - Handles the creation, delivery, and validation of JWT tokens
- Directions
  - Interfaces with the Google Maps API to find the distance of a route
- Journey
  - Provides information about a route including the cost and best driver.
- Roster
  - Handles the store of drivers including adding to roster, removing from roster, and updating price/km

## Docker

The application has been dockerized. 

In order to build the services, use the `docker-compose build` command in the root `easy-ride` directory.

Then, to run the services, use the `docker-compose up` command in the root `easy-ride` directory. 

Note that the `Directions` service requires a Google Maps API key to be set as an environment variable. The easiest way to do this is to add a file `.env` within the `Directions` directory. Within `.env`, set the API key in the format `MAPS_API_KEY=cAbfJkBfABfNAXfaqQvPugjljVV-AquTzpzT1k0`. This is just an example key, you will need to set your own. 