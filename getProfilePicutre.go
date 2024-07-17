package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
)

func HttpToByte(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch URL, status code: %d", response.StatusCode)
	}

	imageBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	return imageBytes
}

func GetProfilePicture(s *Worker, username string) (map[string]interface{}, error) {
	url := "https://www.instagram.com/" + username
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatalf("Failed to fetch page, status code: %d", response.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	var profilePictureURL string
	doc.Find("meta[property='og:image']").Each(func(i int, s *goquery.Selection) {
		content, exists := s.Attr("content")
		if exists {
			profilePictureURL = content
		}
	})

	if profilePictureURL != "" {
		if isExists(s, username) {
			coll := s.client.Database("InstaPFP").Collection("pfp")
			coll.DeleteOne(context.TODO(), bson.D{{Key: "username", Value: username}})
			fmt.Printf("%v Already exists, updating...", username)
		}
		dict := make(map[string]interface{})
		dict["username"] = username
		dict["url"] = profilePictureURL
		dict["byte"] = HttpToByte(profilePictureURL)
		return dict, nil

	} else {
		return nil, errors.New("profile Picture URL not found")
	}
}
