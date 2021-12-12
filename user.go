package main

import (
	"fmt"
  	"gorm.io/gorm"
	"net/http"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Note []Note `json:"-"`
}

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

// Create a struct to read the username and password from the request body
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Register(w http.ResponseWriter, r *http.Request) {
	Setup()

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
	}

	// user.Password = hash(user.Password)
	db.Create(&user)
	fmt.Fprintf(w, "New User successfully registered")
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	Setup()

	var users []User
	db.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func Login(w http.ResponseWriter, r *http.Request) {
	Setup()
	
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User
	db.Where("username = ?", creds.Username).Find(&user)
	expectedPassword := user.Password

	if expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}