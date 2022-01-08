package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/models"

	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/cart"
	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	conn *redis.Client
}

func NewSessionRedisRepository(conn *redis.Client) cart.Repository {
	return &RedisRepository{
		conn: conn,
	}
}

func (r *RedisRepository) getNewKey(value uint64) string {
	return fmt.Sprintf("cart:%d", value)
}

// Delete user cart
func (r *RedisRepository) DeleteCart(userId *models.DtoUserId) error {
	key := r.getNewKey(userId.Id)

	err := r.conn.Del(context.Background(), key).Err()
	if err != nil {
		return errors.ErrDBInternalError
	}
	return nil
}

// Select user cart by id
func (r *RedisRepository) SelectCartById(userId *models.DtoUserId) (*models.DtoCart, error) {
	userCart := &models.DtoCart{}
	key := r.getNewKey(userId.Id)

	data, err := r.conn.Get(context.Background(), key).Bytes()
	if err != nil || data == nil {
		return nil, errors.ErrCartNotFound
	}

	if err = json.Unmarshal(data, userCart); err != nil {
		return nil, errors.ErrCanNotUnmarshal
	}

	return userCart, err
}

// Add new user cart
func (r *RedisRepository) AddCart(userId *models.DtoUserId, userCart *models.DtoCart) error {
	key := r.getNewKey(userId.Id)

	data, err := json.Marshal(userCart)
	if err != nil {
		return errors.ErrCanNotMarshal
	}

	err = r.conn.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		return errors.ErrDBInternalError
	}
	return nil
}
