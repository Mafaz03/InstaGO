package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// type

func getusernameFromHeader(header http.Header) (string, error) {
	val := header.Get("username")
	if val == "" {
		return "", fmt.Errorf("no username found in header, contained: %v", header)
	}
	return val, nil
}

func showPic(w http.ResponseWriter, r *http.Request) {
	username, err := getusernameFromHeader(r.Header)
	if err != nil {
		// fmt.Println("Error: ", err)
		log.Fatal(err)
	}
	pic, err := GetProfilePicture(username)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	dictionary := make(map[string][]byte)
	dictionary[username] = pic
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dictionary)
}

func main() {
	PORT := 8080
	fmt.Println("Listening on PORT: ", PORT)
	http.HandleFunc("/pfp", showPic)
	address := fmt.Sprintf(":%d", PORT)
	http.ListenAndServe(address, nil)
}
