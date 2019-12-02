package main

import (
	"fmt"
	"net/http"
)

var count int
var messages []Message

func Counter(w http.ResponseWriter, r *http.Request) {
	count++
	fmt.Fprintf(w, "%d", count)
}

type Message struct {
	Receiver string
	Text     string
}

func MakeMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create new message")
	text := r.URL.Query().Get("text")
	receiver := r.URL.Query().Get("receiver")

	msg := Message{Receiver: receiver, Text: text}
	messages = append(messages, msg)

}

func main() {
	http.HandleFunc("/Counter", Counter) //endPoint

	http.HandleFunc("/MakeMessage", MakeMessage)
	http.HandleFunc("/PrintMessage", PrintMessage)

	fmt.Println("Server starts at :8080")
	http.ListenAndServe(":8080", nil)
}

// MakeMessage?text=Privet&receiver=Pasha

func PrintMessage(w http.ResponseWriter, r *http.Request) {
	receiver := r.URL.Query().Get("receiver")
	for _, msg := range messages {
		if receiver != "" && receiver != msg.Receiver {
			continue
		}
		fmt.Fprintf(w, "Message \"%s\" sent to  %s\n", msg.Text, msg.Receiver)
	}
}

// pull-request
// залить свой код на репозиторий, гит комит, гит пул
