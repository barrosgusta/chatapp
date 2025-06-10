package main

import (
	"log"
	"net/http"

	"chat-message-service/dynamodb"
	"chat-message-service/handler"
	"chat-message-service/sqs"
)

// corsMiddleware adds CORS headers to the response
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    dynamodb.Init()

    consumer := sqs.NewSQSConsumer()
    go consumer.StartConsuming()

    // mux allows us to use middleware
    mux := http.NewServeMux()
    
    // not used in this version, but can be implemented if needed
    // mux.HandleFunc("/message", handler.HandlePostMessage)

    mux.HandleFunc("/history", handler.HandleGetMessages)

    // Wrap mux with CORS middleware
    handlerWithCORS := corsMiddleware(mux)

    log.Println("Chat Message Service running on :8081")
    if err := http.ListenAndServe(":8081", handlerWithCORS); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}