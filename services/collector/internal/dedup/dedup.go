package dedup

//дедупликация вакансий через Redis
import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Deduplicator struct {
	rdb *redis.Client
	ttl time.Duration // как долго помним что вакансия уже обработана
}

func New(rdb *redis.Client) *Deduplicator {
	return &Deduplicator{
		rdb: rdb,
		ttl: 24 * time.Hour, // помним 24 часа
	}
}

// fingerprint создаёт уникальный ключ для вакансии
// используем source + external_id — этого достаточно для уникальности
func fingerprint(source, externalID string) string {
	h := sha256.Sum256([]byte(source + ":" + externalID))
	return fmt.Sprintf("dedup:%x", h)
}

// IsDuplicate возвращает true если вакансия уже обрабатывалась
// Принцип: пробуем записать ключ в Redis через SetNX (set if not exists)
// Если ключ уже был — значит дубль
func (d *Deduplicator) IsDuplicate(ctx context.Context, source, externalID string) (bool, error) {
	key := fingerprint(source, externalID)

	// SetNX возвращает true если ключ был СОЗДАН (то есть его не было)
	created, err := d.rdb.SetNX(ctx, key, 1, d.ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx: %w", err)
	}

	return !created, nil
}

func (d *Deduplicator) Reset(ctx context.Context, source, externalID string) error {
	key := fingerprint(source, externalID)
	return d.rdb.Del(ctx, key).Err()
}
