package domain

import (
	"fmt"
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewSession(id, name, avatarURL string) (*Session, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("name cannot be empty")
	}

	now := time.Now()
	return &Session{
		ID:        id,
		Name:      name,
		AvatarURL: avatarURL,
		CreatedAt: now,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
	}, nil
}
