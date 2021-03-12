# Easy-Ride

EasyRide is a fictional UK based ride-sharing app powered by microservices written in Go and deployed using Docker.

## API Documentation

Full API documentation can be found in the `API Specification` folder. The API is described by `.yml` folders that have been exported to `html` using Redoc. 

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