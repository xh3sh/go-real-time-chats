package repo

import (
	"context"
	"fmt"
	"time"
)

const Z_MESSAGES_INTERVAL = 30 * time.Minute

func (r *RedisRepository) StartCleanupZMessages(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(Z_MESSAGES_INTERVAL)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_ = r.CleanOldIndices(ctx, Z_MESSAGES_INTERVAL)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (r *RedisRepository) CleanOldIndices(ctx context.Context, retentionPeriod time.Duration) error {
	threshold := time.Now().Add(-retentionPeriod).Unix()

	return r.client.ZRemRangeByScore(
		ctx,
		r.makeKey(ChatMessagesKey),
		"-inf",
		fmt.Sprintf("%d", threshold),
	).Err()
}
