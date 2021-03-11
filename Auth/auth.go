package main

import (
	// "encoding/json"
	"log"
	"net/http"

	//"os"
	"time"

	//"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_jwt_key")

type User struct {
	Username string `json:"username"`
	Name string `json:"name"`
	PasswordHash string
}

var accounts = map[string] User{}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func signIn(w http.ResponseWriter, r *http.Request) {
	// Get given username and password values from request
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Printf("Attempting sign-in of user %s", username)

	// Lookup user in accounts map. Fetch the hashed password
	user := accounts[username]
	hashedPassword := user.PasswordHash

	// If the password does not match the hash, return 401.
	if !verifyPassword(hashedPassword, password) {
		log.Printf("Sign-in of user %s failed.", username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Calculate an expiration time 5 minutes from now
	expirationTime := time.Now().Add(5 * time.Minute)

	// Create a claims struct that includes the username and expiration time. 
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create a token with the HS256 hash method and the claims created above
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		log.Println("Error: Could not create JWT.")
		w.WriteHeader(http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: expirationTime,
	})
}

func hashSaltPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	// If there is an error, pass it on.
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil;
}

func initialiseAccounts() {
	// Controlled case should not have error, safe to ignore here.
	password1, _ := hashSaltPassword("iamjohndoe")
	password2, _ := hashSaltPassword("iambobross")

	accounts["johndoe"] = User {
		Username: "johndoe",
		Name: "John Doe",
		PasswordHash: password1,
	}

	accounts["bobross"] = User {
		Username: "bobross",
		Name: "Bob Ross",
		PasswordHash: password2,
	}

}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", signIn).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	log.Println("Starting Auth Service")
	initialiseAccounts()
	handleRequests()
}