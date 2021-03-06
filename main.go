package main

import (
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
  	"gorm.io/gorm"
)

func Setup() {
	dsn := "host=localhost user=postgres password=postgres dbname=usernotes port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	} else {
		fmt.Println("Successfully connected to database")
	}
}

func InitialMigration() {
	Setup()

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Note{})
}

func HandleRequests() {
	r := mux.NewRouter()

	r.HandleFunc("/register", Register).Methods("POST")
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/users", GetAllUsers).Methods("GET")
	
	r.HandleFunc("/notes", GetNotes).Methods("GET")
	r.HandleFunc("/notes", CreateNote).Methods("POST")
	r.HandleFunc("/notes/{nid}", DeleteNote).Methods("DELETE")
	r.HandleFunc("/notes/{nid}/{title}/{content}", UpdateNote).Methods("PUT")

	fmt.Println("Starting server at 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func main() {
	InitialMigration()
	HandleRequests()
}