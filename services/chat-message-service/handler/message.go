package handler

import (
	"encoding/json"
	"net/http"

	"chat-message-service/dynamodb"
)

// Not used in this version, but can be implemented if needed
// func HandlePostMessage(w http.ResponseWriter, r *http.Request) {
// 	var msg dynamodb.ChatMessage
// 	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
// 		http.Error(w, "invalid request", http.StatusBadRequest)
// 		return
// 	}
// 	msg.ID = uuid.New().String()
// 	msg.Timestamp = time.Now().Format(time.RFC3339)

// 	if err := dynamodb.SaveMessage(msg); err != nil {
// 		http.Error(w, "failed to save message", http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusCreated)
// }

func HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := dynamodb.GetRecentMessages()
	if err != nil {
		http.Error(w, "failed to get messages", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}