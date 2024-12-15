package repository

import (
	"database/sql"
	"time"

	"github.com/SergeyMilch/user-email-verification/internal/models"
	"github.com/google/uuid"
)

type TokenRepository interface {
    CreateForUser(userID string) (*models.EmailVerificationToken, error)
    FindByToken(token string) (*models.EmailVerificationToken, error)
    MarkUsed(tokenID string) error
}

type tokenRepository struct {
    db *sql.DB
}

func NewTokenRepository(db *sql.DB) TokenRepository {
    return &tokenRepository{db: db}
}

func (r *tokenRepository) CreateForUser(userID string) (*models.EmailVerificationToken, error) {
    evToken := &models.EmailVerificationToken{
        ID:        uuid.New().String(),
        UserID:    userID,
        Token:     uuid.New().String(),
        ExpiresAt: time.Now().Add(24 * time.Hour),
        CreatedAt: time.Now().UTC(),
    }

    _, err := r.db.Exec(`INSERT INTO email_verification_tokens (id, user_id, token, expires_at, created_at)
        VALUES ($1, $2, $3, $4, $5)`,
        evToken.ID, evToken.UserID, evToken.Token, evToken.ExpiresAt, evToken.CreatedAt)
    if err != nil {
        return nil, err
    }

    return evToken, nil
}

func (r *tokenRepository) FindByToken(token string) (*models.EmailVerificationToken, error) {
    var ev models.EmailVerificationToken
    var usedAt sql.NullTime
    err := r.db.QueryRow(`SELECT id, user_id, token, expires_at, created_at, used_at
      FROM email_verification_tokens WHERE token=$1`, token).
        Scan(&ev.ID, &ev.UserID, &ev.Token, &ev.ExpiresAt, &ev.CreatedAt, &usedAt)
    if err != nil {
        return nil, err
    }
    if usedAt.Valid {
        ev.UsedAt = &usedAt.Time
    }
    return &ev, nil
}

func (r *tokenRepository) MarkUsed(tokenID string) error {
    now := time.Now().UTC()
    _, err := r.db.Exec(`UPDATE email_verification_tokens SET used_at=$1 WHERE id=$2`, now, tokenID)
    return err
}