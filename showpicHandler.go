package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s *Worker) showPic(w http.ResponseWriter, r *http.Request) {
	username, err := getusernameFromHeader(r.Header)
	if err != nil {
		log.Fatal(err)
	}
	UrlandByte, err := GetProfilePicture(s, username)
	if UrlandByte != nil {
		if err != nil {
			fmt.Println("Error: ", err)
		}
		// dictionary := make(map[string]interface{})
		// dictionary["Data"] = UrlandByte
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UrlandByte)
		SavetoMongo(s, UrlandByte)
	} else {
		fmt.Printf("Could not fetch for %v\n", username)
	}
}
