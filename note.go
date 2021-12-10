package main

import (
	"fmt"
	// "gorm.io/driver/postgres"
  	"gorm.io/gorm"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"strconv"
)

var db *gorm.DB
var err error

type Note struct {
	gorm.Model 
	Title string `gorm:"not null" json:"title"`
	Content string `gorm:"not null" json:"content"`
	UserID uint `gorm:"primaryKey" json:"user_id"`
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

	vars := mux.Vars(r)
	var notes []Note
	db.Where("user_id = ?", vars["uid"]).Find(&notes)

	json.NewEncoder(w).Encode(notes)
}

func CreateNote(w http.ResponseWriter, r *http.Request) { 
	Setup()

	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)

	var note Note
	_ = decoder.Decode(&note)

	val,_ := strconv.Atoi(vars["uid"])
	note.UserID = uint(val)

	db.Create(&note)
	fmt.Fprintf(w, "New notes successfully added")
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	Setup()
	
	vars := mux.Vars(r)

	var note Note
	db.Where("user_id = ? AND id = ?", vars["uid"], vars["nid"]).Find(&note)
	db.Delete(&note)

	fmt.Fprintf(w, "Note successfully deleted")
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	Setup()
	
	vars := mux.Vars(r)
	var note Note

	db.Where("user_id = ? AND id = ?", vars["uid"], vars["nid"]).Find(&note)
	
	note.Title = vars["title"]
	note.Content = vars["content"]

	db.Save(&note)
	fmt.Fprintf(w, "Note successfully updated")
}