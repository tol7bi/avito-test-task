package repository

import (
	"context"
	"errors"
	"pvz-backend/internal/models"
)

func CreateUser(ctx context.Context, email, hash, role string) error {
	_, err := models.DB.Exec(ctx, `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
	`, email, hash, role)
	if err != nil {
		return err
	}
	return nil
}

func GetUserHashAndRole(ctx context.Context, email string) (string, string, error) {
	var hash, role string
	err := models.DB.QueryRow(ctx, `
		SELECT password_hash, role FROM users WHERE email = $1
	`, email).Scan(&hash, &role)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}
	return hash, role, nil
}
