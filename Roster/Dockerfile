FROM golang:1.15

WORKDIR /app/
COPY Roster ./Roster
RUN go get github.com/gorilla/mux 

EXPOSE 8000
CMD ["go", "run", "/app/Roster/roster.go"]

