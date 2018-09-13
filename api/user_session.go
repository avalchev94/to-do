package api

import (
	"time"

	"github.com/google/uuid"
)

func newLoginSession(id int64) (string, error) {
	uuid := uuid.New().String()
	if err := redisConn.Set(uuid, id, 24*time.Hour).Err(); err != nil {
		return "", err
	}
	redisConn.Save()

	return uuid, nil
}

func getSessionUser(uuid string) (int64, error) {
	result := redisConn.Get(uuid)
	if result.Err() != nil {
		return -1, result.Err()
	}
	id, err := result.Int64()
	if err != nil {
		return -1, err
	}

	return id, nil
}
