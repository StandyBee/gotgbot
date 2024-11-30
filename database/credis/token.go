package credis

import (
	"context"
	"fmt"
)

type TokenRepository struct {
	db *RedisDB
}

func NewTokenRepository(db *RedisDB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) SaveAccessToken(chatId int64, token string) error {
	tokenKey := r.tokenKey(chatId, "access_token")
	return r.db.Client.Set(context.Background(), tokenKey, token, 0).Err()
}

func (r *TokenRepository) GetAccessToken(chatId int64) (string, error) {
	return r.db.Client.Get(context.Background(), r.tokenKey(chatId, "access_token")).Result()
}

func (r *TokenRepository) SaveRequestToken(chatId int64, token string) error {
	tokenKey := r.tokenKey(chatId, "request_token")
	return r.db.Client.Set(context.Background(), tokenKey, token, 0).Err()
}

func (r *TokenRepository) GetRequestToken(chatId int64) (string, error) {
	return r.db.Client.Get(context.Background(), r.tokenKey(chatId, "request_token")).Result()
}

func (r *TokenRepository) tokenKey(chatId int64, tokenType string) string {
	return fmt.Sprintf("%d:%s", chatId, tokenType)
}
