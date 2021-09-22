package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type QuestionType struct {
	ID       string `json:ID`
	Question string `json:Question`
	Answer   string `json:Answer`
	Answered bool   `json:Answered`
	User     string `json:User`
}

var questionsPool []QuestionType

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Println("Endpoint Get Hit. Searched ID: ", key)

	for _, question := range questionsPool {
		if question.ID == key {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(question)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "Question not found"}`))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println("Endpoint GetAll Hit")
	json.NewEncoder(w).Encode(questionsPool)
}

func getAllByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	user := vars["user"]

	var userQuestions []QuestionType

	fmt.Println("Endpoint GetAllByUser Hit. Searched user: ", user)

	for _, question := range questionsPool {
		if question.User == user {
			userQuestions = append(userQuestions, question)
		}
	}
	if len(userQuestions) != 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userQuestions)
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "User not found"}`))
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Println("Endpoint create Hit")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var question QuestionType
	json.Unmarshal(reqBody, &question)
	fmt.Println("Storing question", question)
	exists := false
	for _, el := range questionsPool {
		if el.ID == question.ID {
			exists = true
			w.Write([]byte(`{"message": "Unable to create question.  ID already exists!"}`))
		}
	}

	if !exists {
		questionsPool = append(questionsPool, QuestionType{ID: question.ID, Question: question.Question, Answer: "", Answered: false, User: question.User})
		fmt.Println("Creating new question")
	}

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

	for i, q := range questionsPool {
		if q.ID == question.ID {
			questionsPool[i].Question = question.Question
			questionsPool[i].Answer = question.Answer
			if question.Answer == "" {
				questionsPool[i].Answered = false
			} else {
				questionsPool[i].Answered = true
			}
			if question.User != "" {
				questionsPool[i].User = question.User
			}
			fmt.Println("Updating question", question)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Question updated succesfully"}`))
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"message": "Question not found!"}`))
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Println("Endpoint Delete Hit. Requested ID: ", key)

	for i, question := range questionsPool {
		if question.ID == key {
			questionsPool = append(questionsPool[:i], questionsPool[(i+1):]...)
			fmt.Println("Deleting question", question)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Question succesfully deleted."}`))
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"message": "Question not found!"}`))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/qa/").Subrouter()
	api.HandleFunc("/get/{id}", get).Methods(http.MethodGet)
	api.HandleFunc("/getAll", getAll).Methods(http.MethodGet)
	api.HandleFunc("/getAllByUser/{user}", getAllByUser).Methods(http.MethodGet)
	api.HandleFunc("/create", create).Methods(http.MethodPost)
	api.HandleFunc("/update/{id}", update).Methods(http.MethodPut)
	api.HandleFunc("/delete/{id}", delete).Methods(http.MethodDelete)
	api.HandleFunc("", notFound)
	log.Fatal(http.ListenAndServe(":8080", api))

}
