package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
		return 0, errors.New("invalid max header value")
	}
	return maxInt, nil
}

func getrootAccountFromHeader(header http.Header) (string, error) {
	rootAcc := header.Get("root")
	if rootAcc == "" {
		return "", errors.New("invalid root account header")
	}
	return rootAcc, nil
}

func (s *Worker) showPicBatch(w http.ResponseWriter, r *http.Request) {
	rootAcc, err := getrootAccountFromHeader(r.Header)
	if err != nil {
		log.Fatal(err)
	}
	max, err := getmaxFromHeader(r.Header)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("input.txt", []byte(rootAcc), 0644)
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("python3", "getFollowers.py", strconv.Itoa(max))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))

	content, err := ioutil.ReadFile("scraped.txt")
	if err != nil {
		log.Fatal(err)
	}
	usernames := strings.Split(string(content), "\n")

	var success int
	ticker := time.NewTicker(2 * time.Second)
	for num, username := range usernames {
		UrlandByte, err := GetProfilePicture(s, username)
		if UrlandByte != nil {
			if err != nil {
				fmt.Println("Error: ", err)
			}
			SavetoMongo(s, UrlandByte)
			success += 1
		} else {
			fmt.Printf("Could not fetch for %v\n", username)
		}
		if num == max {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "application/json")
			message := map[string]string{
				"Success": fmt.Sprintf("Saved, %v", success),
			}
			json.NewEncoder(w).Encode(message)
			break
		}
		<-ticker.C
	}

}
