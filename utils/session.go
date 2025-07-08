package utils

import (
	"encoding/json"
	"errors"
	"gin-blog-api/database"
	"gin-blog-api/models"
	"time"
)

func SaveSession(token string, userID uint, duration time.Duration) error {
	session := models.Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(duration),
	}
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return database.Rdb.Set(database.Ctx, token, data, duration).Err()
}

func GetSession(token string) (*models.Session, error) {
	val, err := database.Rdb.Get(database.Ctx, token).Result()
	if err != nil {
		return nil, err
	}

	var session models.Session
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return nil, err
	}

	// Süresi dolmuş mu kontrol et
	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session expired")
	}

	return &session, nil
}

func DeleteSession(token string) error {
	return database.Rdb.Del(database.Ctx, token).Err()
}
