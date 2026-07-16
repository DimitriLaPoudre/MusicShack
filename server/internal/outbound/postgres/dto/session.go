package dto

import (
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}

func (s Session) ToSession() model.Session {
	return model.Session{
		ID:        s.ID,
		UserID:    s.UserID,
		Token:     s.Token,
		ExpiresAt: s.ExpiresAt,
	}
}

func SessionsToSessions(dto []Session) []model.Session {
	sessions := []model.Session{}
	for _, s := range dto {
		sessions = append(sessions, model.Session{
			ID:        s.ID,
			UserID:    s.UserID,
			Token:     s.Token,
			ExpiresAt: s.ExpiresAt,
		})
	}
	return sessions
}
