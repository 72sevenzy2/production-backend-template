package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	utils "github.com/72sevenzy2/production-backend-template/internal"
)

type Session struct {
	ExpiresAt time.Time
	UserID    int
}

type Sessions struct {
	mu       sync.RWMutex
	sessions map[string]Session
}

func NewSessions() *Sessions {
	return &Sessions{
		sessions: make(map[string]Session),
	}
}

func (s *Sessions) Create(userID int, sessionID string, expiry time.Time) error {
	if userID <= 0 || sessionID == "" {
		return utils.ErrInvalidSessionData
	}

	if !expiry.After(time.Now()) {
		return utils.ErrInvalidSessionExpiry
	}

	key := hashSID(sessionID)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[key] = Session{
		UserID:    userID,
		ExpiresAt: expiry,
	}

	return nil
}

func (s *Sessions) Get(sessionID string) (Session, error) {
	if sessionID == "" {
		return Session{}, utils.ErrSessionNotFound
	}

	key := hashSID(sessionID)

	s.mu.RLock()
	session, ok := s.sessions[key]
	s.mu.RUnlock()

	if !ok {
		return Session{}, utils.ErrSessionNotFound
	}

	if !session.ExpiresAt.After(time.Now()) {
		// lazy eviction
		s.Delete(sessionID)
		return Session{}, utils.ErrSessionExpired
	}

	return session, nil
}

func (s *Sessions) Delete(sessionID string) error {
	if sessionID == "" {
		return utils.ErrSessionNotFound
	}

	key := hashSID(sessionID)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[key]; !ok {
		return utils.ErrSessionNotFound
	}

	delete(s.sessions, key)
	return nil
}

func hashSID(sessionID string) string {
	sum := sha256.Sum256([]byte(sessionID))
	return hex.EncodeToString(sum[:])
}
