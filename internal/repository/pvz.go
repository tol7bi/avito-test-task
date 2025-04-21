package repository

import (
	"context"
	"fmt"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
	"pvz-backend/internal/models"
)

var allowedCities = map[string]bool{
	"Москва":            true,
	"Санкт-Петербург":   true,
	"Казань":            true,
}

func CreatePVZ(ctx context.Context, db *pgxpool.Pool, city string) (*models.PVZ, error) {
	if !allowedCities[city] {
		return nil, fmt.Errorf("город %s не поддерживается", city)
	}

	var pvz models.PVZ
	err := models.DB.QueryRow(ctx, `
        INSERT INTO pvz (city)
        VALUES ($1)
        RETURNING id, registration_date, city
    `, city).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)

	if err != nil {
		return nil, err
	}
	return &pvz, nil
}


func GetPVZList(ctx context.Context, start, end *time.Time, page, limit int) ([]models.PVZResponse, error) {
	offset := (page - 1) * limit

	var s, e any
	if start != nil {
		s = *start
	}
	if end != nil {
		e = *end
	}

	rows, err := models.DB.Query(ctx, `
	SELECT DISTINCT pvz.id, pvz.registration_date, pvz.city
	FROM pvz
	JOIN receptions r ON r.pvz_id = pvz.id
	WHERE ($1::timestamptz IS NULL OR r.date_time >= $1::timestamptz)
	  AND ($2::timestamptz IS NULL OR r.date_time <= $2::timestamptz)
	ORDER BY pvz.registration_date DESC
	OFFSET $3 LIMIT $4
`, s, e, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.PVZResponse

	for rows.Next() {
		var pvz models.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			return nil, err
		}

		response := models.PVZResponse{PVZ: pvz}

		rx, err := models.DB.Query(ctx, `
			SELECT id, date_time, pvz_id, status
			FROM receptions
			WHERE pvz_id = $1
			ORDER BY date_time DESC
		`, pvz.ID)
		if err != nil {
			return nil, err
		}

		for rx.Next() {
			var r models.Reception
			if err := rx.Scan(&r.ID, &r.DateTime, &r.PVZID, &r.Status); err != nil {
				return nil, err
			}

			prods := []models.Product{}
			prx, err := models.DB.Query(ctx, `
				SELECT id, date_time, type, reception_id
				FROM products
				WHERE reception_id = $1
			`, r.ID)
			if err != nil {
				return nil, err
			}

			for prx.Next() {
				var p models.Product
				if err := prx.Scan(&p.ID, &p.DateTime, &p.Type, &p.ReceptionID); err != nil {
					return nil, err
				}
				prods = append(prods, p)
			}
			prx.Close()

			response.Receptions = append(response.Receptions, models.ReceptionWithProducts{
				Reception: r,
				Products:  prods,
			})
		}
		rx.Close()

		result = append(result, response)
	}

	return result, nil
}

