package service

import (
	"errors"
	"time"

	"github.com/SergeyMilch/user-email-verification/internal/models"
	"github.com/SergeyMilch/user-email-verification/internal/repository"
)

type UserService interface {
    RegisterUser(nick, name, email string) (*models.User, string, error)
    VerifyEmail(token string) (*models.User, error)
}

type userService struct {
    userRepo  repository.UserRepository
    tokenRepo repository.TokenRepository
}

func NewUserService(u repository.UserRepository, t repository.TokenRepository) UserService {
    return &userService{
        userRepo:  u,
        tokenRepo: t,
    }
}

func (s *userService) RegisterUser(nick, name, email string) (*models.User, string, error) {
    user := &models.User{
        Nick:          nick,
        Name:          name,
        Email:         email,
        EmailVerified: false,
    }

    err := s.userRepo.Create(user)
    if err != nil {
        return nil, "", err
    }

    evToken, err := s.tokenRepo.CreateForUser(user.ID)
    if err != nil {
        return nil, "", err
    }

	// Здесь можно было бы отправить настоящее письмо
    // Пока просто возвращаем ссылку в ответе

    return user, evToken.Token, nil
}

func (s *userService) VerifyEmail(token string) (*models.User, error) {
    evToken, err := s.tokenRepo.FindByToken(token)
    if err != nil {
        return nil, errors.New("invalid token")
    }

    if time.Now().After(evToken.ExpiresAt) {
        return nil, errors.New("token expired")
    }

    if evToken.UsedAt != nil {
        return nil, errors.New("token already used")
    }

    user, err := s.userRepo.FindByID(evToken.UserID)
    if err != nil {
        return nil, errors.New("user not found")
    }

    if user.EmailVerified {
        return nil, errors.New("email already verified")
    }

    user.EmailVerified = true
    err = s.userRepo.Update(user)
    if err != nil {
        return nil, errors.New("failed to update user")
    }

    err = s.tokenRepo.MarkUsed(evToken.ID)
    if err != nil {
        return nil, errors.New("failed to mark token as used")
    }

    return user, nil
}