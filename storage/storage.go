package storage

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/max2sax/raft-chat/models"
)

type Storage struct {
	rooms    sync.Map // map[string]*models.Room
	messages sync.Map // map[string][]models.Message
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) CreateRoom(name, description string) *models.Room {
	id := uuid.New().String()
	room := &models.Room{
		ID:          id,
		Name:        name,
		Description: description,
	}
	s.rooms.Store(id, room)
	s.messages.Store(id, []models.Message{})
	return room
}

func (s *Storage) AddMessage(roomID, sender, content string) error {
	_, ok := s.rooms.Load(roomID)
	if !ok {
		return fmt.Errorf("room not found")
	}
	msg := models.Message{
		Timestamp: time.Now(),
		Sender:    sender,
		Content:   content,
	}
	val, _ := s.messages.Load(roomID)
	msgs := val.([]models.Message)
	msgs = append(msgs, msg)
	s.messages.Store(roomID, msgs)
	return nil
}

func (s *Storage) GetMessages(roomID string) ([]models.Message, error) {
	_, ok := s.rooms.Load(roomID)
	if !ok {
		return nil, fmt.Errorf("room not found")
	}
	val, _ := s.messages.Load(roomID)
	msgs := val.([]models.Message)
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
