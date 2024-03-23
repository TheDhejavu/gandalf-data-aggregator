package store

import (
	"fmt"
	"sync"
)

type SessionStore struct {
	sessions map[string]string
	mutex    sync.RWMutex
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]string),
	}
}

// GetSession retrieves session data by session ID.
func (store *SessionStore) GetSession(sessionID string) (string, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	sessionData, ok := store.sessions[sessionID]
	if !ok {
		return "", fmt.Errorf("session not found")
	}
	return sessionData, nil
}

// SaveSession saves session data with the given session ID.
func (store *SessionStore) SaveSession(sessionID, sessionData string) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.sessions[sessionID] = sessionData
	return nil
}

func (store *SessionStore) FindKeyOrSaveSession(sessionID, sessionData string) (string, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, ok := store.sessions[sessionID]; ok {
		return sessionID, nil
	}

	for id, data := range store.sessions {
		if data == sessionData {
			return id, nil
		}
	}

	store.sessions[sessionID] = sessionData
	return sessionID, nil
}
