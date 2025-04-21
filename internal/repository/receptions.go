package repository

import (
	"context"
	"errors"
	"pvz-backend/internal/models"
)

func CreateReception(ctx context.Context, pvzID string) (*models.Reception, error) {
	var exists bool
	err := models.DB.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM receptions
            WHERE pvz_id = $1 AND status = 'in_progress'
        )
    `, pvzID).Scan(&exists)

	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("нельзя создать новую приёмку: предыдущая не закрыта")
	}

	var r models.Reception
	err = models.DB.QueryRow(ctx, `
        INSERT INTO receptions (pvz_id, status)
        VALUES ($1, 'in_progress')
        RETURNING id, date_time, pvz_id, status
    `, pvzID).Scan(&r.ID, &r.DateTime, &r.PVZID, &r.Status)

	if err != nil {
		return nil, err
	}
	return &r, nil
}

func CloseLastReception(ctx context.Context, pvzID string) (*models.Reception, error) {
	var r models.Reception

	err := models.DB.QueryRow(ctx, `
		SELECT id, date_time, pvz_id, status
		FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC
		LIMIT 1
	`, pvzID).Scan(&r.ID, &r.DateTime, &r.PVZID, &r.Status)

	if err != nil {
		return nil, errors.New("нет активной приёмки для закрытия")
	}

	_, err = models.DB.Exec(ctx, `
		UPDATE receptions
		SET status = 'close'
		WHERE id = $1
	`, r.ID)

	if err != nil {
		return nil, err
	}

	r.Status = "close"
	return &r, nil
}
