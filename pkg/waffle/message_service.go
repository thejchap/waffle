package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// MessageCapacity Total capacity of the data store,
// chosen in tandem with MessageIDPossibilities to provide a reasonable storage
// size (4096 most recent messages) with a low probability of ID collisions
// from client-generated IDs
const MessageCapacity int = 4096

// MessageIDPossibilities Number of possible message IDs, based on frontend ID
// generation algorithm
const MessageIDPossibilities = 68719476736

// MessageService Stateful service containing messages data store
type MessageService struct {
	messages []Message
	onCreate func(Message)
	capacity int
}

// Create and return a new Service with the store initialized with a maximum
// capacity
func newMessageService(onCreate func(Message)) *MessageService {
	service := &MessageService{
		messages: make([]Message, 0, MessageCapacity),
		onCreate: onCreate,
		capacity: MessageCapacity,
	}

	prob := collisionProbability(service)

	log.Print("[waffle/message_service] Initialized")
	log.Printf("[waffle/message_service] Store capacity: %d messages", MessageCapacity)
	log.Printf("[waffle/message_service] ID collision probability: %f", prob)

	return service
}

// Handler for a POST request containing message attributes. Return the newly
// created message
func (s *MessageService) createHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var msg Message
	err := decoder.Decode(&msg)

	if err != nil {
		panic(err)
	}

	// Truncate the list, keep to only the most recent {MessageCapacity} messages
	if s.isFull() {
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
func collisionProbability(service *MessageService) float32 {
	var result float32 = 1.0

	for i := 0; i < service.capacity; i++ {
		result *= (1 - float32(i)/MessageIDPossibilities)
	}

	return 1.0 - result
}

// Return whether or not the data store is full
func (s *MessageService) isFull() bool {
	return len(s.messages) == s.capacity
}
