// go
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// --- Mock SQSProducer ---
type mockSQSProducer struct {
	mu      sync.Mutex
	messages []string
}

func (m *mockSQSProducer) SendMessage(ctx context.Context, msg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, msg)
	return nil
}

// --- Helper functions ---
func startTestServer(hub *Hub) *httptest.Server {
	r := http.NewServeMux()
	r.HandleFunc("/ws", wsHandler(hub))
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return httptest.NewServer(r)
}

func wsConnect(t *testing.T, url string) *websocket.Conn {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("WebSocket dial failed: %v", err)
	}
	return conn
}

// --- Tests ---

func TestHealthzEndpoint(t *testing.T) {
	hub := NewHub()
	ts := startTestServer(hub)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/healthz")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserRegistrationAndBroadcast(t *testing.T) {
	hub := NewHub()
	hub.sqsProducer = &mockSQSProducer{}
	go hub.Run()
	ts := startTestServer(hub)
	defer ts.Close()

	wsURL := "ws" + ts.URL[4:] + "/ws"
	conn := wsConnect(t, wsURL)
	defer conn.Close()

	// Set name
	   setName := map[string]any{"type": "SET_NAME", "name": "alice"}
	   assert.NoError(t, conn.WriteJSON(setName))

	// Expect NAME_ACCEPTED
	var resp map[string]any
	assert.NoError(t, conn.ReadJSON(&resp))
	t.Logf("Received: %#v", resp)
	if resp["type"] == nil {
		t.Fatalf("Expected message with type, got: %#v", resp)
	}
	assert.Equal(t, "NAME_ACCEPTED", resp["type"])

	// Expect USER_LIST broadcast
	assert.NoError(t, conn.ReadJSON(&resp))
	t.Logf("Received: %#v", resp)
	if resp["type"] == nil {
		t.Fatalf("Expected message with type, got: %#v", resp)
	}
	assert.Equal(t, "USER_LIST", resp["type"])
	assert.Contains(t, resp["users"], "alice")
}

// func TestNameRejected(t *testing.T) {
// 	hub := NewHub()
// 	hub.sqsProducer = &mockSQSProducer{}
// 	go hub.Run()
// 	ts := startTestServer(hub)
// 	defer ts.Close()

// 	wsURL := "ws" + ts.URL[4:] + "/ws"
// 	conn1 := wsConnect(t, wsURL)
// 	defer conn1.Close()
// 	conn2 := wsConnect(t, wsURL)
// 	defer conn2.Close()

// 	// First client sets name
// 	   conn1.WriteJSON(map[string]any{"type": "SET_NAME", "name": "bob"})
// 	var resp map[string]any
// 	conn1.ReadJSON(&resp)
// 	conn1.ReadJSON(&resp) // USER_LIST

// 	// Second client tries same name
// 	   conn2.WriteJSON(map[string]any{"type": "SET_NAME", "name": "bob"})
// 	// Loop to find NAME_REJECTED or fail after a few tries
// 	for i := 0; i < 5; i++ {
// 		conn2.ReadJSON(&resp)
// 		t.Logf("Received: %#v", resp)
// 		if resp["type"] == "NAME_REJECTED" {
// 			break
// 		}
// 		if i == 4 {
// 			t.Fatalf("Did not receive NAME_REJECTED after 5 messages, last: %#v", resp)
// 		}
// 	}
// 	assert.Equal(t, "NAME_REJECTED", resp["type"])
// }

func TestMessageBroadcastAndSQS(t *testing.T) {
	mockSQS := &mockSQSProducer{}
	hub := NewHub()
	hub.sqsProducer = mockSQS
	go hub.Run()
	ts := startTestServer(hub)
	defer ts.Close()

	wsURL := "ws" + ts.URL[4:] + "/ws"
	conn := wsConnect(t, wsURL)
	defer conn.Close()

	   conn.WriteJSON(map[string]any{"type": "SET_NAME", "name": "carol"})
	conn.ReadJSON(&map[string]any{}) // NAME_ACCEPTED
	conn.ReadJSON(&map[string]any{}) // USER_LIST

	// Send message
	msgText := "hello world"
	   conn.WriteJSON(map[string]any{"type": "MESSAGE", "text": msgText})

	// Expect MESSAGE broadcast
	var resp map[string]any
	assert.NoError(t, conn.ReadJSON(&resp))
	t.Logf("Received: %#v", resp)
	if resp["type"] == nil {
		t.Fatalf("Expected message with type, got: %#v", resp)
	}
	assert.Equal(t, "MESSAGE", resp["type"])
	msgMap, ok := resp["message"].(map[string]any)
	if !ok {
		t.Fatalf("Expected message to be map[string]interface{}, got: %#v", resp["message"])
	}
	assert.Equal(t, "carol", msgMap["user"])
	assert.Equal(t, msgText, msgMap["text"])

	// Wait for SQS goroutine
	time.Sleep(100 * time.Millisecond)
	mockSQS.mu.Lock()
	defer mockSQS.mu.Unlock()
	assert.Len(t, mockSQS.messages, 1)
	t.Logf("SQS message: %s", mockSQS.messages[0])
	var sqsMsg map[string]any
	assert.NoError(t, json.Unmarshal([]byte(mockSQS.messages[0]), &sqsMsg))
	t.Logf("Unmarshaled SQS message: %#v", sqsMsg)
	// Accept both "text" and "Text" keys for robustness
	if val, ok := sqsMsg["text"]; ok {
		assert.Equal(t, msgText, val)
	} else if val, ok := sqsMsg["Text"]; ok {
		assert.Equal(t, msgText, val)
	} else {
		t.Fatalf("SQS message missing 'text' or 'Text' key: %#v", sqsMsg)
	}
}

func TestTypingBroadcast(t *testing.T) {
	hub := NewHub()
	hub.sqsProducer = &mockSQSProducer{}
	go hub.Run()
	ts := startTestServer(hub)
	defer ts.Close()

	wsURL := "ws" + ts.URL[4:] + "/ws"
	conn := wsConnect(t, wsURL)
	defer conn.Close()

	conn.WriteJSON(map[string]any{"type": "SET_NAME", "name": "dave"})
	conn.ReadJSON(&map[string]any{}) // NAME_ACCEPTED
	conn.ReadJSON(&map[string]any{}) // USER_LIST

	// Send TYPING_START
	conn.WriteJSON(map[string]any{"type": "TYPING_START"})
	var resp map[string]any
	assert.NoError(t, conn.ReadJSON(&resp))
	assert.Equal(t, "TYPING", resp["type"])
	t.Logf("Typing users after start: %#v", resp["typing"])
	assert.Contains(t, resp["typing"], "dave")
}