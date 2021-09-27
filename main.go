package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QuestionType struct {
	ID       string `json:ID`
	Question string `json:Question`
	Answer   string `json:Answer`
	Answered bool   `json:Answered`
	User     string `json:User`
}

var collection *mongo.Collection

//var clientOptions, client, connection, collection, err interface{}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Println("Endpoint Get Hit. Searched ID: ", key)

	filter := bson.D{{"id", key}}

	var result QuestionType
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println("Question not found!")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Question not found"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Question found:", result)
	json.NewEncoder(w).Encode(result)

}

func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Endpoint GetAll Hit")
	findOptions := options.Find()

	var results []*QuestionType
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		fmt.Println("Error with database!")
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(`{"message": "Error with database!"}`))
		return
	}

	for cur.Next(context.TODO()) {
		var elem QuestionType
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if len(results) == 0 {
		fmt.Println("Empty database!")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Empty database!"}`))
		return
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)

}

func getAllByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	user := vars["user"]

	fmt.Println("Endpoint GetAllByUser Hit. Searched user: ", user)

	findOptions := options.Find()
	filter := bson.D{{"user", user}}
	var results []*QuestionType
	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		fmt.Println("User not found!")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "User not found!"}`))
		return
	}

	for cur.Next(context.TODO()) {
		var elem QuestionType
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if len(results) == 0 {
		fmt.Println("User not found!!")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "User not found!"}`))
		return
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)

}

func create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint create Hit")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var question QuestionType
	json.Unmarshal(reqBody, &question)
	fmt.Println("Requested question", question)

	filter := bson.D{{"id", question.ID}}

	var result QuestionType
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == nil {
		fmt.Printf("ID already exists!: %+v\n", result)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Unable to create question.  ID already exists!"}`))
		return
	}

	_, err = collection.InsertOne(context.TODO(), question)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Unable to create question.  Error!"}`))
		return
	}

	fmt.Println("Storing question")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Question posted succesfully"}`))

}

func update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Println("Endpoint Update Hit. Requested ID: ", key)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var question QuestionType
	json.Unmarshal(reqBody, &question)

	if question.ID != key {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Invalid input"}`))
		fmt.Printf("Invalid input. Key %v, ID %v\n", key, question.ID)
		return
	}

	filter := bson.D{{"id", question.ID}}
	fmt.Println("filter", filter)

	if question.Answer != "" {
		question.Answered = true
	} else {
		question.Answered = false
	}

	update := bson.D{{"$set", bson.D{
		{"question", question.Question},
		{"answer", question.Answer},
		{"answered", question.Answered},
		{"user", question.User},
	}}}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "DB Error!"}`))
		return
	}

	if updateResult.MatchedCount == 0 {
		fmt.Printf("ID %v not found!\n", question.ID)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Question not found!"}`))
		return
	}

	fmt.Println("Updating question", question)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Question updated succesfully"}`))

}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Println("Endpoint Delete Hit. Requested ID: ", key)
	filter := bson.D{{"id", key}}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "DB Error!"}`))
		return
	}

	if deleteResult.DeletedCount == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Question not found!")
		w.Write([]byte(`{"message": "Question not found!"}`))
		return
	}

	fmt.Println("Deleting question", key)
	fmt.Printf("Deleted %v document in the collection\n", deleteResult.DeletedCount)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Question succesfully deleted."}`))

}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

func NewRouter() (s *mux.Router) {
	r := mux.NewRouter()
	api := r.PathPrefix("/qa/").Subrouter()
	api.HandleFunc("/get/{id}", get).Methods(http.MethodGet)
	api.HandleFunc("/getAll", getAll).Methods(http.MethodGet)
	api.HandleFunc("/getAllByUser/{user}", getAllByUser).Methods(http.MethodGet)
	api.HandleFunc("/create", create).Methods(http.MethodPost)
	api.HandleFunc("/update/{id}", update).Methods(http.MethodPut)
	api.HandleFunc("/delete/{id}", delete).Methods(http.MethodDelete)
	api.HandleFunc("", notFound)

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://user1:user1!@bd-go-level-vii.oq2dd.mongodb.net/myFirstDatabase")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	// Disconnect defer
	/*defer func() {
		err = client.Disconnect(context.TODO())

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection to MongoDB closed.")
	}()*/

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database("db_questions").Collection("collection_qa")
	fmt.Println("Server Ready")

	return api
}

func main() {

	api := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", api))

}
