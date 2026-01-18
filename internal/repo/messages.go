package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/xh3sh/go-real-time-chats/internal/emmit"
)

const ChatMessagesKey = "chat_messages"

type RedisRepository struct {
	client    *redis.Client
	keyPrefix string
}

func NewRedisRepository(client *redis.Client, keyPrefix string) *RedisRepository {
	return &RedisRepository{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

func (r *RedisRepository) makeKey(key string) string {
	return fmt.Sprintf("%s:%s", r.keyPrefix, key)
}

func (r *RedisRepository) SaveMessage(ctx context.Context, msg emmit.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return r.client.ZAdd(ctx, r.makeKey(ChatMessagesKey), redis.Z{
		Score:  float64(msg.Timestamp),
		Member: data,
	}).Err()
}

func (r *RedisRepository) GetMessages(ctx context.Context) ([]emmit.Message, error) {
	data, err := r.client.ZRange(ctx, r.makeKey(ChatMessagesKey), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messages []emmit.Message
	for _, item := range data {
		var msg emmit.Message
		if err := json.Unmarshal([]byte(item), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *RedisRepository) GetMessagesAfter(ctx context.Context, timestamp int64) ([]emmit.Message, error) {
	data, err := r.client.ZRangeByScore(ctx, r.makeKey(ChatMessagesKey), &redis.ZRangeBy{
		Min: fmt.Sprintf("(%d", timestamp),
		Max: "+inf",
	}).Result()

	if err != nil {
		return nil, err
	}

	var messages []emmit.Message
	for _, item := range data {
		var msg emmit.Message
		if err := json.Unmarshal([]byte(item), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
