package main

import (
	"fmt"
  	"gorm.io/gorm"
	"net/http"
	"encoding/json"
)



type User struct {
	gorm.Model
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Note []Note `json:"-"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	Setup()

	decoder := json.NewDecoder(r.Body)
	
	var user User
	_ = decoder.Decode(&user)

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
	
	decoder := json.NewDecoder(r.Body)

	var user User
	_ = decoder.Decode(&user)

	var dbUser User
	db.Where("username = ?", user.Username).Find(&dbUser)

	// if dbUser.Password == hash(user.Password) {
	// 	// jwt := GenerateJWT(dbUser.ID)  
	// } else {
	// 	fmt.Fprintf(w, "Invalid Password")
	// }
}