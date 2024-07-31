package storage

import (
	"sync"
	"user-service/internal/user"
)

type Storage struct {
	mu    sync.Mutex
	users map[string]user.User
}

func New() *Storage {
	return &Storage{
		users: make(map[string]user.User),
	}
}

func (s *Storage) AddUser(u user.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.ID] = u
}

func (s *Storage) GetUser(id string) (user.User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, exists := s.users[id]
	return u, exists
}

func (s *Storage) UpdateUser(u user.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.ID] = u
}

func (s *Storage) DeleteUser(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.users, id)
}

func (s *Storage) ListUsers() []user.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	users := make([]user.User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users
}
