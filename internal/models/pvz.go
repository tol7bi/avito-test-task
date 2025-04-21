package models

import "time"

type CreatePVZRequest struct {
	City string `json:"city"`
}

type PVZ struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

type PVZResponse struct {
	PVZ        PVZ                    `json:"pvz"`
	Receptions []ReceptionWithProducts `json:"receptions"`
}

type ReceptionWithProducts struct {
	Reception Reception `json:"reception"`
	Products  []Product `json:"products"`
}
