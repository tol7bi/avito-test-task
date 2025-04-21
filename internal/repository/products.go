package repository

import (
	"context"
	"errors"
	"pvz-backend/internal/models"
)

func AddProduct(ctx context.Context, pvzID string, prodType string) (*models.Product, error) {
	// незакрытая приёмку
	var receptionID string
	err := models.DB.QueryRow(ctx, `
		SELECT id FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC
		LIMIT 1
	`, pvzID).Scan(&receptionID)

	if err != nil {
		return nil, errors.New("нет активной приёмки для этого ПВЗ")
	}

	var product models.Product
	err = models.DB.QueryRow(ctx, `
		INSERT INTO products (reception_id, type)
		VALUES ($1, $2)
		RETURNING id, date_time, type, reception_id
	`, receptionID, prodType).Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionID)

	if err != nil {
		return nil, err
	}
	return &product, nil
}


func DeleteLastProduct(ctx context.Context, pvzID string) error {
	var receptionID string
	err := models.DB.QueryRow(ctx, `
		SELECT id FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC
		LIMIT 1
	`, pvzID).Scan(&receptionID)

	if err != nil {
		return errors.New("нет активной приёмки")
	}

	var productID string
	err = models.DB.QueryRow(ctx, `
		SELECT id FROM products
		WHERE reception_id = $1
		ORDER BY date_time DESC
		LIMIT 1
	`, receptionID).Scan(&productID)

	if err != nil {
		return errors.New("нет товаров для удаления")
	}

	_, err = models.DB.Exec(ctx, `DELETE FROM products WHERE id = $1`, productID)
	if err != nil {
		return err
	}

	return nil
}