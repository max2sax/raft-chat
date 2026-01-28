package storage

import (
	"fmt"
	"sort"
	"sync"

	"github.com/max2sax/raft-chat/models"
	"github.com/oklog/ulid/v2"
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
		_, ok := s.rooms.Load(req.message.RoomID)
		if !ok {
			req.result <- fmt.Errorf("room not found")
			continue
		}

		val, _ := s.messages.Load(req.message.RoomID)
		msgs := val.([]models.Message)
		msgs = append(msgs, *req.message)
		s.messages.Store(req.message.RoomID, msgs)

		req.result <- nil
	}
}

func (s *Storage) CreateRoom(name, description string) *models.Room {
	id := ulid.Make().String()
	room := &models.Room{
		ID:          id,
		Name:        name,
		Description: description,
	}
	s.rooms.Store(id, room)
	s.messages.Store(id, []models.Message{})
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

func (s *Storage) GetMessages(roomID string) ([]models.Message, error) {
	_, ok := s.rooms.Load(roomID)
	if !ok {
		return nil, fmt.Errorf("room not found")
	}
	val, _ := s.messages.Load(roomID)
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
