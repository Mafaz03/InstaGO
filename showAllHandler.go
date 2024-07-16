package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Worker) showAll(w http.ResponseWriter, r *http.Request) {
	coll := s.client.Database("InstaPFP").Collection("pfp")
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Println("Database is empty")
	}
	results := []bson.M{}
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Println("Unable to fetch")
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)

}
