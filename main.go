package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type User struct {
	UserID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name"`
	Email  string             `json:"email" bson:"email"`
	Passwd string             `json:"password,omitempty" bson:"hash_password"`
}
type Post struct {
	PostID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserIDPost primitive.ObjectID `json:"userpostid" bson:"userpostid"`
	Image      string             `json:"image" bson:"image"`
	Caption    string             `json:"caption" bson:"caption"`
	TimePosted time.Time          `json:"posttime,omitempty" bson:"posttime,omitempty"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	coll := client.Database("AppointyInsta").Collection("Users")
	res, err := coll.InsertOne(context.TODO(), user)

	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Document inserted with ID: %s\n", res.InsertedID)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	coll := client.Database("AppointyInsta").Collection("Posts")
	post.TimePosted = time.Now().UTC()
	res, err := coll.InsertOne(context.TODO(), post)

	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Document inserted with ID: %s\n", res.InsertedID)
}

func getUserID(w http.ResponseWriter, r *http.Request) {
	var user User

	coll := client.Database("AppointyInsta").Collection("Users")

	id := r.URL.Path[len("/users/"):]
	objectID, err := primitive.ObjectIDFromHex(id)

	err = coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectID}}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", objectID)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
}

func getPostID(w http.ResponseWriter, r *http.Request) {
	var post Post

	id := r.URL.Path[len("/posts/"):]
	coll := client.Database("AppointyInsta").Collection("Posts")

	objectID, err := primitive.ObjectIDFromHex(id)

	err = coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectID}}).Decode(&post)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", objectID)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(post, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
}

func listUserPosts(w http.ResponseWriter, r *http.Request) {
	//list all the posts of user
}
func handleRequests() {
	http.HandleFunc("/users", createUser)
	http.HandleFunc("/users/", getUserID)
	http.HandleFunc("/posts", createPost)
	http.HandleFunc("/posts/", getPostID)
	http.HandleFunc("/posts/users/", listUserPosts)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
func main() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/")
	}
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
	handleRequests()
}
