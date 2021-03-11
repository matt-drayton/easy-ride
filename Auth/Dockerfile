FROM golang:1.15

WORKDIR /app/
COPY Auth ./Auth
RUN go get github.com/gorilla/mux golang.org/x/crypto/bcrypt github.com/dgrijalva/jwt-go

EXPOSE 8000
CMD ["go", "run", "/app/Auth/auth.go"]

