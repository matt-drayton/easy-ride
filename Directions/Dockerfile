FROM golang:1.15

WORKDIR /app/
COPY Directions ./Directions
RUN go get github.com/gorilla/mux github.com/kr/pretty googlemaps.github.io/maps

EXPOSE 8000
CMD ["go", "run", "/app/Directions/directions.go"]