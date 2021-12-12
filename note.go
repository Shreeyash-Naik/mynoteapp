package main

import (
	"fmt"
  	"gorm.io/gorm"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"os"
)

var db *gorm.DB
var err error

type Note struct {
	gorm.Model 
	Title string `gorm:"unique not null" json:"title"`
	Content string `gorm:"not null" json:"content"`
	UserID uint `gorm:"primaryKey" json:"user_id"`
}

func Authorize(w http.ResponseWriter, r *http.Request) (string) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			os.Exit(1)
		}
		w.WriteHeader(http.StatusBadRequest)
		os.Exit(1)
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			os.Exit(1)
		}
		w.WriteHeader(http.StatusBadRequest)
		os.Exit(1)
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
	}

	return claims.Username
}

func GetNotes(w http.ResponseWriter, r *http.Request) {
	Setup()

	username := Authorize(w, r)

	var user User
	db.Where("username = ?", username).Find(&user)

	var notes []Note
	db.Where("user_id = ?", user.ID).Find(&notes)

	fmt.Fprintf(w, "Hello %s\n", username)
	json.NewEncoder(w).Encode(notes)
}

func CreateNote(w http.ResponseWriter, r *http.Request) { 
	Setup()

	username := Authorize(w, r)

	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)

	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
	}

	var user User
	db.Where("username = ?", username).Find(&user)

	note.UserID = user.ID

	db.Create(&note)
	fmt.Fprintf(w, "New notes successfully added")
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	Setup()
	
	username := Authorize(w, r)

	var user User
	db.Where("username = ?", username).Find(&user)

	vars := mux.Vars(r)
	var note Note
	db.Where("user_id = ? AND id = ?", user.ID, vars["nid"]).Find(&note)
	db.Delete(&note)

	fmt.Fprintf(w, "Note successfully deleted")
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	Setup()
	username := Authorize(w, r)

	var user User
	db.Where("username = ?", username).Find(&user)

	vars := mux.Vars(r)
	var note Note
	db.Where("user_id = ? AND id = ?", user.ID, vars["nid"]).Find(&note)
	
	note.Title = vars["title"]
	note.Content = vars["content"]

	db.Save(&note)
	fmt.Fprintf(w, "Note successfully updated")
}