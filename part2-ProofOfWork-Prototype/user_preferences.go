package main

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

type UserPreferences struct {
    UserID    string `json:"user_id"`
    Theme     string `json:"theme"`
    Language  string `json:"language"`
    Notifications bool `json:"notifications"`
}

type PreferenceCache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewPreferenceCache(addr string, ttl time.Duration) *PreferenceCache {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: "",
        DB:       0,
    })
    return &PreferenceCache{client: rdb, ttl: ttl}
}

func (c *PreferenceCache) Get(ctx context.Context, userID string) (*UserPreferences, error) {
    key := fmt.Sprintf("prefs:%s", userID)
    val, err := c.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    } else if err != nil {
        return nil, err
    }

    var prefs UserPreferences
    if err := json.Unmarshal([]byte(val), &prefs); err != nil {
        return nil, err
    }
    return &prefs, nil
}

func (c *PreferenceCache) Set(ctx context.Context, prefs *UserPreferences) error {
    key := fmt.Sprintf("prefs:%s", prefs.UserID)
    data, err := json.Marshal(prefs)
    if err != nil {
        return err
    }
    return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *PreferenceCache) Invalidate(ctx context.Context, userID string) error {
    key := fmt.Sprintf("prefs:%s", userID)
    return c.client.Del(ctx, key).Err()
}