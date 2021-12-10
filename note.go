package main

import (
	"fmt"
	// "gorm.io/driver/postgres"
  	"gorm.io/gorm"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
)

var db *gorm.DB
var err error

type Note struct {
	gorm.Model 
	Title string `gorm:"not null" json:"title"`
	Content string `gorm:"not null" json:"content"`
}

// func (n *Note) UnmarshalJSON(data []byte) error {
//     var s [2]string
//     if err := json.Unmarshal(data, &s); err != nil {
//         return err
//     }

//     n.Title = s[0]
//     n.Content = s[1]
//     return nil
// }

func GetNotes(w http.ResponseWriter, r *http.Request) {
	Setup()

	// params := mux.Vars(r)
	var notes []Note
	db.Find(&notes)

	json.NewEncoder(w).Encode(notes)
}

func CreateNote(w http.ResponseWriter, r *http.Request) { 
	Setup()

	decoder := json.NewDecoder(r.Body)
	var n Note
	_ = decoder.Decode(&n)

	db.Create(&n)
	fmt.Fprintf(w, "New notes successfully added")
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	Setup()
	
	vars := mux.Vars(r)

	var note Note
	db.Where("id = ?", vars["nid"]).Find(&note)
	db.Delete(&note)

	fmt.Fprintf(w, "Note successfully deleted")
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	Setup()
	
	vars := mux.Vars(r)

	var note Note
	// Find the note using note id
	db.Where("id = ?", vars["nid"]).Find(&note)
	
	note.Title = vars["title"]
	note.Content = vars["content"]

	db.Save(&note)
	fmt.Fprintf(w, "Note successfully updated")
}