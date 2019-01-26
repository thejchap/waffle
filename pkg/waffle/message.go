package server

// Message main message data type
type Message struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
