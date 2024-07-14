package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func GetProfilePicture(username string) (string, error) {
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
		return "", err
	}

	var profilePictureURL string
	doc.Find("meta[property='og:image']").Each(func(i int, s *goquery.Selection) {
		content, exists := s.Attr("content")
		if exists {
			profilePictureURL = content
		}
	})

	if profilePictureURL != "" {
		// fmt.Println("Profile Picture URL:", profilePictureURL)
		return profilePictureURL, nil
	} else {
		// fmt.Println("Profile Picture URL not found")
		return "", errors.New("profile Picture URL not found")
	}

}
