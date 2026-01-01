package request

import "github.com/google/uuid"

type AkunGrupCreateRequest struct {
	Nama string `validate:"required,max=64" json:"nama"`
	Kode string `validate:"required,max=12" json:"kode"`
}

type AkunGrupUpdateRequest struct {
	Id   uuid.UUID
	Nama string `validate:"required,max=64" json:"nama"`
	Kode string `validate:"required,max=12" json:"kode"`
}
