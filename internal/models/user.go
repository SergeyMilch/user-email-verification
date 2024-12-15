package models

import "time"

type User struct {
    ID            string    `json:"id"`
    Nick          string    `json:"nick"`
    Name          string    `json:"name"`
    Email         string    `json:"email"`
    EmailVerified bool      `json:"email_verified"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}