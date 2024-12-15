package models

import "time"

type EmailVerificationToken struct {
    ID        string     `json:"id"`
    UserID    string     `json:"user_id"`
    Token     string     `json:"token"`
    ExpiresAt time.Time  `json:"expires_at"`
    CreatedAt time.Time  `json:"created_at"`
    UsedAt    *time.Time `json:"used_at,omitempty"`
}