package network

import (
	"context"
	"time"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	redis "github.com/go-redis/redis/v8"
)

type Redis struct {
	URL    string
	logger logging.Logger
	rdb    *redis.Client
}

var ctx = context.Background()

// NewRedis creates a new Redis instances and accepts a URL encoded string to configure the
// connection with Redis service.
// redis://<user>:<pass>@localhost:6379/<db>
func NewRedis(url string, logger logging.Logger) *Redis {
	return &Redis{url, logger, nil}
}

// Start starts a connection with the Redis service by parsing the configuration URL and creating
// a new client instance, which is added to the struct responsible for abstracting this service
// capabilities.
func (r *Redis) Start(started chan bool) {
	opt, err := redis.ParseURL(r.URL)
	if err != nil {
		r.logger.Error(err)
	}

	r.rdb = redis.NewClient(opt)
	started <- true
}

// Set stores a key-value pair to the Redis database. The key must be a string and the value can be
// of any supported type: https://redis.io/topics/data-types.
func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	return r.rdb.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from the Redis database according to key, which is returned as a string.
func (r *Redis) Get(key string) (string, error) {
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}

	return val, nil
}
