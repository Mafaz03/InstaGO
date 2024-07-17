package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Worker struct {
	client *mongo.Client
}

func NewWorker(c *mongo.Client) *Worker {
	return &Worker{
		client: c,
	}
}

func SavetoMongo(s *Worker, info map[string]interface{}) {
	coll := s.client.Database("InstaPFP").Collection("pfp")
	jsonData, err := json.Marshal(info)
	if err != nil {
		log.Fatalf("Error marshalling info: %v", err)
	}
	var Info bson.M
	err = json.Unmarshal(jsonData, &Info)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON to BSON: %v", err)
	}
	_, err = coll.InsertOne(context.TODO(), Info)
	if err != nil {
		log.Fatalf("Error Could not InsertOne: %v", err)
	}
	fmt.Printf("%v\nSuccessfully Saved\n\n", info["username"])
}

func getusernameFromHeader(header http.Header) (string, error) {
	val := header.Get("username")
	if val == "" {
		return "", fmt.Errorf("no username found in header, contained: %v", header)
	}
	return val, nil
}

type Format struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Byte     []byte             `bson:"byte,omitempty"`
	URL      string             `bson:"url,omitempty"`
	Username string             `bson:"username,omitempty"`
}

func isExists(s *Worker, username string) bool {
	coll := s.client.Database("InstaPFP").Collection("pfp")

	var result Format
	err := coll.FindOne(context.TODO(), bson.D{{Key: "username", Value: username}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		return false
	}
	return true
}

func main() {
	godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("Port not found in env")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Could not connect to Mongo Server")
	}
	server := NewWorker(client)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	
	router.Get("/pfp", server.showPic)
	router.Get("/showall", server.showAll)
	router.Get("/showallBatches", server.showPicBatch)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}
	log.Printf("listening on Port number: %v", PORT)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	

	// http.HandleFunc("/pfp", server.showPic)
	// http.HandleFunc("/showall", server.showAll)
	// address := fmt.Sprintf(":%v", PORT)
	// http.ListenAndServe(address, nil)
}
