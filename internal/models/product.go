package models

import "time"

type Product struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionID string    `json:"receptionId"`
}

type CreateProductRequest struct {
	Type  string `json:"type"`
	PVZID string `json:"pvzId"`
}
