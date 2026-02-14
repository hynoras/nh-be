package infra

import (
	"context"
	"nh-be/utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionStore interface {
	CreateUserSession(ctx context.Context, sessionId string, userId uuid.UUID) error
	GetUserSession(ctx context.Context, sessionId string) (uuid.UUID, error)
	DeleteUserSession(ctx context.Context, sessionId string) error
}

type sessionStore struct {
	rdb *redis.Client
}

func NewSessionStore(rdb *redis.Client) SessionStore {
	return &sessionStore{rdb: rdb}
}

func (s *sessionStore) CreateUserSession(ctx context.Context, sessionId string, userId uuid.UUID) error {
	rdbErr := s.rdb.Set(ctx, "session:"+sessionId, userId.String(), 8*time.Hour).Err()
	if rdbErr != nil {
		return rdbErr
	}
	return nil
}

func (s *sessionStore) GetUserSession(ctx context.Context, sessionId string) (uuid.UUID, error) {
	userId, err := s.rdb.Get(ctx, "session:"+sessionId).Result()
	if err != nil {
		return uuid.Nil, err
	}

	parsedUuid, err := utils.ParseStringToUUID(userId)
	if err != nil {
		return uuid.Nil, err
	}
	return parsedUuid, nil
}

func (s *sessionStore) DeleteUserSession(ctx context.Context, sessionId string) error {
	err := s.rdb.Del(ctx, "session:"+sessionId).Err()
	if err != nil {
		return err
	}
	return nil
}
