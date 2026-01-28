package models

import "time"

type Room struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
}
