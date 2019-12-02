package main

import (
	"fmt"
	"net/http"
	"strconv"
)

var messages []Message

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

	http.HandleFunc("/MakeMessage", MakeMessage)
	http.HandleFunc("/PrintMessage", PrintMessage)
	http.HandleFunc("/AutoMessage", AutoMessage)

	fmt.Println("Server starts at :8080")
	http.ListenAndServe(":8080", nil)
}
func PrintMessage(w http.ResponseWriter, r *http.Request) {
	receiver := r.URL.Query().Get("receiver")
	for _, msg := range messages {
		if receiver != "" && receiver != msg.Receiver {
			continue
		}
		fmt.Fprintf(w, "Message \"%s\" sent to  %s\n", msg.Text, msg.Receiver)
	}
}
func AutoMessage(w http.ResponseWriter, r *http.Request) {
	// localhost:8080/AutoMessage?text=Privet&receiver=Pasha&quantity=5
	quantity, err := strconv.Atoi(r.URL.Query().Get("quantity"))
	if err != nil {
		fmt.Fprintf(w, "Incorrect enter of quantity! Quantity is a number")
	} else {
		text := r.URL.Query().Get("text")
		receiver := r.URL.Query().Get("receiver")
		msg := Message{Receiver: receiver, Text: text}
		fmt.Fprintf(w, "Message is \"%s\", receiver is \"%s\", quantity of repeated messages is %d", text, receiver, quantity)
		for i := 0; i < quantity; i++ {
			messages = append(messages, msg)
		}
	}
}

// Branch testing get a comment
