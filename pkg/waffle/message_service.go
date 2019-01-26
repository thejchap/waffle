package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// Total capacity of the data store, chosen in tandem with MessageIDPossibilities to
// provide a reasonable storage size (4096 most recent messages) with a low
// probability of ID collisions from client-generated IDs
const MessageCapacity = 4096

// Number of possible message IDs, based on frontend ID generation algorithm
const MessageIDPossibilities = 68719476736

// Stateful service
type MessageService struct {
	messages []Message
	onCreate func(Message)
}

// Create and return a new Service with the store initialized with a maximum
// capacity
func newMessageService(onCreate func(Message)) *MessageService {
	service := &MessageService{
		messages: make([]Message, 0, MessageCapacity),
		onCreate: onCreate,
	}

	prob := collisionProbability()

	log.Print("[waffle/message_service] Initialized")
	log.Printf("[waffle/message_service] Store capacity: %d messages", MessageCapacity)
	log.Printf("[waffle/message_service] ID collision probability: %f", prob)

	return service
}

// Handler for a POST request containing message attributes. Return the newly
// created message
func (s *MessageService) createHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	msg := Message{}
	err := decoder.Decode(&msg)

	if err != nil {
		panic(err)
	}

	// Truncate the list, keep to only the most recent {MessageCapacity} messages
	if len(s.messages) == MessageCapacity {
		s.messages = s.messages[1:len(s.messages)]
	}

	s.messages = append(s.messages, msg)

	// Printing all messages in plain text to logs...Compliance department may
	// not approve :(
	log.Printf("[waffle/message_service] Message created: %#v", msg)

	if s.onCreate != nil {
		s.onCreate(msg)
	}

	json.NewEncoder(w).Encode(msg)
}

// Handler for a GET request. Return all messages in the store
func (s *MessageService) indexHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.messages)
}

// Calculate chance of collisions given the maximum capacity of the data store
// and the number of possible IDs. https://en.wikipedia.org/wiki/Birthday_problem
func collisionProbability() float32 {
	var result float32 = 1.0

	for i := 0; i < MessageCapacity; i++ {
		result *= (1 - float32(i)/MessageIDPossibilities)
	}

	return 1.0 - result
}
