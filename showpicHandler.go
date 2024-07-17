package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func getmaxFromHeader(header http.Header) (int, error) {
	maxStr := header.Get("max")
	maxInt, err := strconv.Atoi(maxStr)
	if err != nil {
		return 0, errors.New("Invalid max header value")
	}
	return maxInt, nil
}

func (s *Worker) showPicBatch(w http.ResponseWriter, r *http.Request) {

	// usernames := []string{"nigg.pablo", "virat.kohli", "anushkasharma"}
	content, err := ioutil.ReadFile("followers_scrape/downloads/scraped.txt")
	if err != nil {
		log.Fatal(err)
	}
	usernames := strings.Split(string(content), "\n")
	max, err := getmaxFromHeader(r.Header)
	if err != nil {
		log.Fatal(err)
	}
	for num, username := range usernames {
		UrlandByte, err := GetProfilePicture(s, username)
		if UrlandByte != nil {
			if err != nil {
				fmt.Println("Error: ", err)
			}
			SavetoMongo(s, UrlandByte)
		} else {
			fmt.Printf("Could not fetch for %v\n", username)
		}
		if num == max {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "application/json")
			message := map[string]string{
				"Success": fmt.Sprintf("Saved, %v", num),
			}
			json.NewEncoder(w).Encode(message)
			break
		}
	}

}
