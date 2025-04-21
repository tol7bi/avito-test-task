package models

import "time"

type CreateReceptionRequest struct {
	PVZID string `json:"pvzId"`
}

type Reception struct {
	ID        string    `json:"id"`
	DateTime  time.Time `json:"dateTime"`
	PVZID     string    `json:"pvzId"`
	Status    string    `json:"status"`
}
