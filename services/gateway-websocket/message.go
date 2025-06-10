package main

type IncomingMessage struct {
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
	Text string `json:"text,omitempty"`
}

type OutgoingMessage struct {
	Type      string      `json:"type"`
	Message   *ChatMessage `json:"message,omitempty"`
	Users     []string    `json:"users,omitempty"`
	Reason    string      `json:"reason,omitempty"`
	Typing    []string    `json:"typing,omitempty"`
}

type ChatMessage struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}
