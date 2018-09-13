package main

import (
	"time"

	"github.com/google/uuid"
)

// NewLoginSession creates session uuid for the user id.
func NewLoginSession(id int64) (string, error) {
	uuid := uuid.New().String()
	if err := RedisConn.Set(uuid, id, 24*time.Hour).Err(); err != nil {
		return "", err
	}
	RedisConn.Save()

	return uuid, nil
}

// GetSessionUser returns the user's id mapped to the uuid.
func GetSessionUser(uuid string) (int64, error) {
	result := RedisConn.Get(uuid)
	if result.Err() != nil {
		return -1, result.Err()
	}
	id, err := result.Int64()
	if err != nil {
		return -1, err
	}

	return id, nil
}
