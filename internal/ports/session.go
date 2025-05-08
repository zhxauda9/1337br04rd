package ports

import (
	"context"
	"net/http"

	"1337b04rd/internal/domain"
)

type SessionService interface {
	CreateSession(ctx context.Context, name string) (*domain.Session, error)
	UpdateSession(w http.ResponseWriter, r *http.Request, newName string) (*domain.Session, error)
	DeleteSession(w http.ResponseWriter, r *http.Request) error
	GetSession(w http.ResponseWriter, r *http.Request) (*domain.Session, error)
	GetAllSessions(ctx context.Context) ([]*domain.Session, error)
}
