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
	mux     *http.ServeMux
	server  *http.Server
}

func NewAPI(store *storage.Storage, server *http.Server) *API {
	mux := http.NewServeMux()
	server.Handler = mux
	return &API{
		storage: store,
		mux:     mux,
		server:  server,
	}
}

func (a *API) RegisterRoutes() {
	a.mux.HandleFunc("POST /rooms", a.createRoomHandler)
	a.mux.HandleFunc("POST /rooms/{roomID}/messages", a.addMessageHandler)
	a.mux.HandleFunc("GET /rooms/{roomID}/messages", a.getMessagesHandler)
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
	roomID := r.PathValue("roomID")
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
	roomID := r.PathValue("roomID")
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
