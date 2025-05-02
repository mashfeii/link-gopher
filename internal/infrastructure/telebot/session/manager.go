package session

import (
	"sync"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
)

type Manager struct {
	sessions map[int64]models.Session
	mu       sync.RWMutex
}

func NewSessionManager() *Manager {
	return &Manager{
		sessions: make(map[int64]models.Session),
	}
}

func (s *Manager) Get(chatID int64) models.Session {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.sessions[chatID]
}

func (s *Manager) GetTrack(chatID int64) (*models.TrackSession, bool) {
	if session := s.Get(chatID); session != nil {
		if session, ok := session.(*models.TrackSession); ok {
			return session, true
		}
	}

	return nil, false
}

func (s *Manager) GetListUntrack(chatID int64) (*models.ListUntrackSession, bool) {
	if session := s.Get(chatID); session != nil {
		if session, ok := session.(*models.ListUntrackSession); ok {
			return session, true
		}
	}

	return nil, false
}

func (s *Manager) Set(chatID int64, session models.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[chatID] = session
}

func (s *Manager) Clear(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, chatID)
}
