package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	Sender User   `json:"sender"`
	Text   string `json:"text"`
}

var messages = make([]Message, 0)

func HandleMessage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
			return
		}
		var msg Message
		if err := json.Unmarshal(body, &msg); err != nil {
			http.Error(w, "Error unmarshaling request body",
				http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		messages = append(messages, msg)
	case "GET":
		jsonBody, err := json.Marshal(messages)
		if err != nil {
			http.Error(w, "Error converting results to json",
				http.StatusInternalServerError)
		}
		w.Write(jsonBody)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

func main() {

	http.HandleFunc("/api/messages", HandleMessage)

	fmt.Println("Server starts at :8080")
	http.ListenAndServe(":8080", nil)
}
