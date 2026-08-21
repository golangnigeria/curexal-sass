package repository

import (
	"context"
	"fmt"

	"github.com/golangnigeria/curexal/internal/kernel/server"
)

type SessionRepository struct {
	server *server.Server
}

func NewSessionRepository(s *server.Server) *SessionRepository {
	return &SessionRepository{server: s}
}

// RevokeOtherUserSessions revokes/deletes all active sessions for a user EXCEPT the current session token or session ID.
func (r *SessionRepository) RevokeOtherUserSessions(ctx context.Context, userID string, activeSessionIdentifier string) error {
	db := r.server.DB.Conn(ctx)
	
	query := `
		DELETE FROM identity.sessions 
		WHERE user_id = $1 
		  AND token != $2 
		  AND id != $2
	`
	_, err := db.Exec(ctx, query, userID, activeSessionIdentifier)
	if err != nil {
		return fmt.Errorf("failed to revoke user sessions: %w", err)
	}
	return nil
}

// RevokeAllUserSessions revokes/deletes ALL active sessions for a user.
func (r *SessionRepository) RevokeAllUserSessions(ctx context.Context, userID string) error {
	db := r.server.DB.Conn(ctx)
	query := `DELETE FROM identity.sessions WHERE user_id = $1`
	_, err := db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all user sessions: %w", err)
	}
	return nil
}
