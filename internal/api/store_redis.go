package api

import (
	"context"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"
)

type ScriptRedisStore struct {
	Group string
	Redis *redis.Client
}

func NewScriptRedisStore(groupName string, redis *redis.Client) *ScriptRedisStore {
	return &ScriptRedisStore{Group: groupName, Redis: redis}
}

func (s *ScriptRedisStore) Load(callback ScriptLoadCallback) {
	scriptNames, err := s.List()
	if err != nil {
		log.Printf("[gojs] Failed to list scripts from Redis: %v", err)
		return
	}

	for _, name := range scriptNames {
		code, err := s.Get(name)
		if err != nil {
			log.Printf("[gojs] Failed to get script '%s' from Redis: %v", name, err)
			continue
		}

		callback(name, code)
	}
}

func (s *ScriptRedisStore) Save(scriptName string, scriptCode string) error {
	if s.Redis == nil {
		return errors.New("redis client not initialized")
	}

	err := s.Redis.HSet(context.Background(), s.Group, scriptName, scriptCode).Err()
	if err != nil {
		log.Printf("[gojs] Failed to store script %s in Redis: %v", scriptName, err)
		return err
	}

	return nil
}

func (s *ScriptRedisStore) Get(scriptName string) (string, error) {
	if s.Redis == nil {
		return "", errors.New("redis client not initialized")
	}

	scriptCode, err := s.Redis.HGet(context.Background(), s.Group, scriptName).Result()
	if err != nil {
		return "", err
	}

	return scriptCode, nil
}

func (s *ScriptRedisStore) Delete(scriptName string) error {
	if s.Redis == nil {
		return errors.New("redis client not initialized")
	}

	err := s.Redis.HDel(context.Background(), s.Group, scriptName).Err()
	if err != nil {
		log.Printf("[gojs] Failed to delete script %s from Redis: %v", scriptName, err)
		return err
	}

	return nil
}

func (s *ScriptRedisStore) List() ([]string, error) {
	if s.Redis == nil {
		return nil, errors.New("redis client not initialized")
	}

	scriptNames, err := s.Redis.HKeys(context.Background(), s.Group).Result()
	if err != nil {
		log.Printf("[gojs] Failed to list scripts from Redis: %v", err)
		return nil, err
	}

	return scriptNames, nil
}

// ScriptExists checks if a script exists in Redis
func (s *ScriptRedisStore) Exists(scriptName string) (bool, error) {
	if s.Redis == nil {
		return false, errors.New("redis client not initialized")
	}

	exists, err := s.Redis.HExists(context.Background(), s.Group, scriptName).Result()
	if err != nil {
		return false, err
	}

	return exists, nil
}
