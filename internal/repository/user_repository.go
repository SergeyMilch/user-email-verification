package repository

import (
	"database/sql"
	"time"

	"github.com/SergeyMilch/user-email-verification/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
    Create(user *models.User) error
    Update(user *models.User) error
    FindByID(id string) (*models.User, error)
}

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
    user.ID = uuid.New().String()
    now := time.Now().UTC()
    user.CreatedAt = now
    user.UpdatedAt = now
    _, err := r.db.Exec(`INSERT INTO users (id, nick, name, email, email_verified, created_at, updated_at)
     VALUES ($1, $2, $3, $4, $5, $6, $7)`,
        user.ID, user.Nick, user.Name, user.Email, user.EmailVerified, user.CreatedAt, user.UpdatedAt)
    return err
}

func (r *userRepository) Update(user *models.User) error {
    user.UpdatedAt = time.Now().UTC()
    _, err := r.db.Exec(`UPDATE users SET nick=$1, name=$2, email=$3, email_verified=$4, updated_at=$5 WHERE id=$6`,
        user.Nick, user.Name, user.Email, user.EmailVerified, user.UpdatedAt, user.ID)
    return err
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
    var u models.User
    err := r.db.QueryRow(`SELECT id, nick, name, email, email_verified, created_at, updated_at FROM users WHERE id=$1`, id).
        Scan(&u.ID, &u.Nick, &u.Name, &u.Email, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return &u, nil
}