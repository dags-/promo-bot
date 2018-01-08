package server

import (
	"sync"
	"time"
)

type AuthSessions struct {
	lock      sync.RWMutex
	timeout   time.Duration
	sessions  map[string]time.Time
	cooldowns map[string]time.Time
}

func newAuthSessions() AuthSessions {
	return AuthSessions{
		timeout:   time.Duration(time.Minute * 30),
		sessions:  make(map[string]time.Time),
		cooldowns: make(map[string]time.Time),
	}
}

func (s *AuthSessions) setAuthenticated(id string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sessions[id] = time.Now()
}

func (s *AuthSessions) isAuthenticated(id string) (bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.sessions[id]
	return ok
}

func (s *AuthSessions) dropAuthentication(id string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.sessions, id)
}

func (s *AuthSessions) setRateLimited(id string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cooldowns[id] = time.Now()
}

func (s *AuthSessions) isRateLimited(id string) (bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if timestamp, ok := s.cooldowns[id]; ok {
		return time.Since(timestamp).Minutes() < s.timeout.Minutes()
	}

	return false
}

func (s *AuthSessions) tick() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for id, timestamp := range s.sessions {
		duration := time.Since(timestamp)
		if duration.Minutes() > s.timeout.Minutes() {
			delete(s.sessions, id)
		}
	}

	for id, timestamp := range s.cooldowns {
		duration := time.Since(timestamp)
		if duration.Minutes() > s.timeout.Minutes() {
			delete(s.cooldowns, id)
		}
	}
}
