package repositories

import (
	"context"
	"go-EdTech/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionsRepository struct {
	db *pgxpool.Pool
}

func NewSessionsRepository(conn *pgxpool.Pool) *SessionsRepository {
	return &SessionsRepository{db: conn}
}

func (r *SessionsRepository) CreateSession(c context.Context, session models.Session) error {
	_, err := r.db.Exec(c,
		`INSERT INTO sessions (user_uuid, refresh_token, expires_at) 
		 VALUES ($1, $2, $3)`,
		session.UserUUID, session.RefreshToken, session.ExpiresAt)
	return err
}

func (r *SessionsRepository) GetSession(c context.Context, refreshToken string) (models.Session, int, error) {
	var session models.Session
	var roleID int

	err := r.db.QueryRow(c,
		`SELECT s.id, s.user_uuid, s.refresh_token, s.expires_at, u.role_id 
         FROM sessions s
         JOIN users u ON s.user_uuid = u.uuid
         WHERE s.refresh_token = $1 AND s.expires_at > NOW()`,
		refreshToken).
		Scan(&session.ID, &session.UserUUID, &session.RefreshToken, &session.ExpiresAt, &roleID)

	return session, roleID, err
}

func (r *SessionsRepository) UpdateSession(c context.Context, session models.Session) error {
	_, err := r.db.Exec(c,
		`UPDATE sessions 
		 SET refresh_token = $1, expires_at = $2 
		 WHERE user_uuid = $3`,
		session.RefreshToken, session.ExpiresAt, session.UserUUID)
	return err
}

func (r *SessionsRepository) DeleteSession(c context.Context, refreshToken string) error {
	_, err := r.db.Exec(c,
		`DELETE FROM sessions 
		 WHERE refresh_token = $1`,
		refreshToken)
	return err
}
