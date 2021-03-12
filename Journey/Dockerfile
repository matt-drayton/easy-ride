FROM golang:1.15

WORKDIR /app/
COPY Journey ./Journey
RUN go get github.com/gorilla/mux 

EXPOSE 8000
CMD ["go", "run", "/app/Journey/journey.go"]