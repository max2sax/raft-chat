package storage

import (
	"fmt"
	"sort"
	"sync"

	"github.com/max2sax/raft-chat/models"
)

type messageWriteRequest struct {
	message *models.Message
	result  chan error
}

type Storage struct {
	rooms            sync.Map // map[string]*models.Room
	messages         sync.Map // map[string][]models.Message
	messageWriteChan chan messageWriteRequest
}

func NewStorage() *Storage {
	s := &Storage{
		messageWriteChan: make(chan messageWriteRequest),
	}
	go s.messageWriter()
	return s
}

func (s *Storage) messageWriter() {
	for req := range s.messageWriteChan {
		_, ok := s.rooms.Load(req.message.RoomName)
		if !ok {
			req.result <- fmt.Errorf("room not found")
			continue
		}

		val, _ := s.messages.Load(req.message.RoomName)
		msgs := val.([]models.Message)
		msgs = append(msgs, *req.message)
		s.messages.Store(req.message.RoomName, msgs)

		req.result <- nil
	}
}

func (s *Storage) CreateRoom(name string, description *string) *models.Room {
	room := &models.Room{
		ID:   name,
		Name: name,
	}
	// Store or update if exists
	rm, loaded := s.rooms.Load(name)
	if loaded {
		if description != nil {
			room.Description = *description
			// Update description if room already exists
			s.rooms.Store(name, room)
		} else {
			room.Description = rm.(*models.Room).Description
		}
		return room
	}
	if description != nil {
		room.Description = *description
	}
	s.rooms.Store(name, room)
	// Initialize messages if room is new
	s.messages.Store(name, []models.Message{})
	return room
}

func (s *Storage) AddMessage(message *models.Message) error {
	result := make(chan error, 1)
	s.messageWriteChan <- messageWriteRequest{
		message: message,
		result:  result,
	}
	return <-result
}

func (s *Storage) GetMessages(roomName string) ([]models.Message, error) {
	_, ok := s.rooms.Load(roomName)
	if !ok {
		return nil, fmt.Errorf("room not found")
	}
	val, _ := s.messages.Load(roomName)
	msgs := val.([]models.Message)
	// Sort by ID ascending (chronological)
	// theoretically these are already sorted by ID except maybe the end
	// , but sorting to be safe
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].ID < msgs[j].ID
	})
	// Return up to 20 messages
	if len(msgs) > 20 {
		return msgs[len(msgs)-20:], nil
	}
	return msgs, nil
}

func (s *Storage) GetRoom(roomName string) (*models.Room, error) {
	room, ok := s.rooms.Load(roomName)
	if !ok {
		return nil, fmt.Errorf("room not found")
	}
	return room.(*models.Room), nil
}

func (s *Storage) GetAllRooms() []*models.Room {
	var rooms []*models.Room
	s.rooms.Range(func(key, value interface{}) bool {
		rooms = append(rooms, value.(*models.Room))
		return true
	})
	return rooms
}
