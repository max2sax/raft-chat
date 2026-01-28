package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

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

type CreateRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CreateMessageRequest struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type Storage struct {
	rooms    sync.Map // map[string]*Room
	messages sync.Map // map[string][]Message
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) CreateRoom(name, description string) *Room {
	id := uuid.New().String()
	room := &Room{
		ID:          id,
		Name:        name,
		Description: description,
	}
	s.rooms.Store(id, room)
	s.messages.Store(id, []Message{})
	return room
}

func (s *Storage) AddMessage(roomID, sender, content string) error {
	_, ok := s.rooms.Load(roomID)
	if !ok {
		return fmt.Errorf("room not found")
	}
	msg := Message{
		Timestamp: time.Now(),
		Sender:    sender,
		Content:   content,
	}
	val, _ := s.messages.Load(roomID)
	msgs := val.([]Message)
	msgs = append(msgs, msg)
	s.messages.Store(roomID, msgs)
	return nil
}

func (s *Storage) GetMessages(roomID string) ([]Message, error) {
	_, ok := s.rooms.Load(roomID)
	if !ok {
		return nil, fmt.Errorf("room not found")
	}
	val, _ := s.messages.Load(roomID)
	msgs := val.([]Message)
	// Sort by timestamp ascending (chronological)
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Timestamp.Before(msgs[j].Timestamp)
	})
	// Return up to 20 messages
	if len(msgs) > 20 {
		return msgs[len(msgs)-20:], nil
	}
	return msgs, nil
}

var storage = NewStorage()

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	room := storage.CreateRoom(req.Name, req.Description)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

func addMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := r.URL.Path[len("/rooms/"):]
	if roomID == "" {
		http.Error(w, "Room ID required", http.StatusBadRequest)
		return
	}

	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Sender == "" || req.Content == "" {
		http.Error(w, "Sender and content are required", http.StatusBadRequest)
		return
	}

	if err := storage.AddMessage(roomID, req.Sender, req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := r.URL.Path[len("/rooms/"):]
	if roomID == "" {
		http.Error(w, "Room ID required", http.StatusBadRequest)
		return
	}

	msgs, err := storage.GetMessages(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}

func main() {
	http.HandleFunc("/rooms", createRoomHandler)
	http.HandleFunc("/rooms/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 7 && path[7:15] == "/messages" { // /rooms/{id}/messages
			roomID := path[7 : len(path)-9] // extract roomID
			r.URL.Path = "/rooms/" + roomID + "/messages"
			if r.Method == http.MethodPost {
				addMessageHandler(w, r)
			} else if r.Method == http.MethodGet {
				getMessagesHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.NotFound(w, r)
		}
	})

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
