package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already in use")
)

type AuthService struct {
	users    *queries.UserQueries
	sessions *queries.SessionQueries
}

func NewAuthService(users *queries.UserQueries, sessions *queries.SessionQueries) *AuthService {
	return &AuthService{users: users, sessions: sessions}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*model.User, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("hash password: %w", err)
	}

	user, err := s.users.Create(ctx, name, email, string(hash))
	if err != nil {
		return nil, "", ErrEmailTaken
	}

	token, err := s.createSession(ctx, user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("create session: %w", err)
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*model.User, string, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.createSession(ctx, user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("create session: %w", err)
	}

	return user, token, nil
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*model.User, error) {
	tokenHash := HashToken(token)

	session, err := s.sessions.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("session not found")
	}

	user, err := s.users.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	tokenHash := HashToken(token)
	session, err := s.sessions.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil
	}
	return s.sessions.Delete(ctx, session.ID)
}

func (s *AuthService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.users.GetPermissions(ctx, userID)
}

func (s *AuthService) createSession(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	tokenHash := HashToken(token)
	_, err = s.sessions.Create(ctx, userID, tokenHash)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
