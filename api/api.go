package api

import (
	"encoding/json"
	"net/http"

	"github.com/max2sax/raft-chat/storage"
)

type CreateRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CreateMessageRequest struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type API struct {
	storage *storage.Storage
	server  *http.Server
}

func NewAPI(store *storage.Storage, server *http.Server) *API {
	return &API{
		storage: store,
		server:  server,
	}
}

func (a *API) RegisterRoutes() {
	http.HandleFunc("/rooms", a.createRoomHandler)
	http.HandleFunc("/rooms/", a.roomsHandler)
}

func (a *API) Start() error {
	return a.server.ListenAndServe()
}

func (a *API) createRoomHandler(w http.ResponseWriter, r *http.Request) {
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

	room := a.storage.CreateRoom(req.Name, req.Description)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

func (a *API) addMessageHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := a.storage.AddMessage(roomID, req.Sender, req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := r.URL.Path[len("/rooms/"):]
	if roomID == "" {
		http.Error(w, "Room ID required", http.StatusBadRequest)
		return
	}

	msgs, err := a.storage.GetMessages(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}

func (a *API) roomsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if len(path) > 7 && path[7:15] == "/messages" {
		roomID := path[7 : len(path)-9]
		r.URL.Path = "/rooms/" + roomID + "/messages"
		if r.Method == http.MethodPost {
			a.addMessageHandler(w, r)
		} else if r.Method == http.MethodGet {
			a.getMessagesHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		http.NotFound(w, r)
	}
}
