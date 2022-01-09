package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/DuckLuckBreakout/web/backend/services/session/pkg/models"
	"github.com/DuckLuckBreakout/web/backend/services/session/pkg/session"
	"github.com/DuckLuckBreakout/web/backend/services/session/server/errors"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	conn *redis.Client
}

func NewSessionRedisRepository(conn *redis.Client) session.Repository {
	return &RedisRepository{
		conn: conn,
	}
}

func (r *RedisRepository) getNewKey(value string) string {
	return fmt.Sprintf("session:%s", value)
}

// Add user session in repository
func (r *RedisRepository) AddSession(session *models.Session) error {
	data := fmt.Sprintf("%d", session.UserId)
	key := r.getNewKey(session.Value)

	err := r.conn.Set(context.Background(), key, data, models.ExpireSessionCookie*time.Second).Err()
	if err != nil {
		return errors.ErrDBInternalError
	}

	return nil
}

// Get user global id from redis db
func (r *RedisRepository) SelectUserIdBySession(sessionValue string) (uint64, error) {
	key := r.getNewKey(sessionValue)

	data, err := r.conn.Get(context.Background(), key).Bytes()
	if err != nil {
		return 0, errors.ErrSessionNotFound
	}

	userId, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return 0, errors.ErrInternalError
	}

	return userId, nil
}

// Delete session from db
func (r *RedisRepository) DeleteSessionByValue(sessionValue string) error {
	key := r.getNewKey(sessionValue)

	err := r.conn.Del(context.Background(), key).Err()
	if err != nil {
		return errors.ErrDBInternalError
	}

	return nil
}
