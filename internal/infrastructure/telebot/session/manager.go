package session

import (
	"sync"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
)

type Manager struct {
	sessions map[int64]*models.UserSession
	mu       sync.Mutex
}

func NewSessionManager() *Manager {
	return &Manager{
		sessions: make(map[int64]*models.UserSession),
	}
}

func (s *Manager) Get(chatID int64) *models.UserSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[chatID]
	if !exists {
		return nil
	}

	return session
}

func (s *Manager) Set(chatID int64, session *models.UserSession) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[chatID] = session
}

func (s *Manager) Clear(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, chatID)
}
