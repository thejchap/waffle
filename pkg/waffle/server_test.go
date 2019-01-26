package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMessageServiceIndexHandler(t *testing.T) {
	service := newMessageService(nil)
	msg := Message{ID: "1", Sender: "1", Content: "Test", Timestamp: 12345}
	service.messages = append(service.messages, msg)

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(service.indexHandler)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("expected status code %v got %v", http.StatusOK, status)
	}

	expected := [1]Message{msg}

	var response [1]Message
	err = json.Unmarshal([]byte(recorder.Body.String()), &response)

	if response != expected {
		t.Errorf(
			"invalid body. expected %v got %v",
			expected,
			recorder.Body.String(),
		)
	}
}

func TestMessageServiceCreateHandler(t *testing.T) {
	service := newMessageService(nil)

	// Test truncation
	service.capacity = 1

	postBody1 := `{"id": "1","sender":"1","content":"Test","timestamp":12345}`
	postBody2 := `{"id": "2","sender":"1","content":"Test","timestamp":12346}`

	req1, err1 := http.NewRequest("POST", "/", strings.NewReader(postBody1))
	req2, err2 := http.NewRequest("POST", "/", strings.NewReader(postBody2))

	if err1 != nil {
		t.Fatal(err1)
	} else if err2 != nil {
		t.Fatal(err2)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(service.createHandler)

	handler.ServeHTTP(recorder, req1)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("expected status code %v got %v", http.StatusOK, status)
	}

	handler.ServeHTTP(recorder, req2)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("expected status code %v got %v", http.StatusOK, status)
	}

	if len(service.messages) > 1 {
		t.Error("expected old message to be truncated")
	}

	msg := Message{ID: "2", Sender: "1", Content: "Test", Timestamp: 12346}

	if service.messages[0] != msg {
		t.Error("expected message to be created")
	}
}

// TODO: Write tests for SSEBroker
func TestSSEBroker(t *testing.T) {}
