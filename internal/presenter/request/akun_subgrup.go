package request

import "github.com/google/uuid"

type AkunSubgrupCreateRequest struct {
	IdAkunGrup uuid.UUID `validate:"required" json:"id_akun_grup"`
	Kode       string    `validate:"required,min=1,max=255" json:"kode"`
	Nama       string    `validate:"required,min=1,max=255" json:"nama"`
}

type AkunSubgrupUpdateRequest struct {
	Id         uuid.UUID
	IdAkunGrup uuid.UUID `validate:"required" json:"id_akun_grup"`
	Kode       string    `validate:"required,min=1,max=255" json:"kode"`
	Nama       string    `validate:"required,min=1,max=255" json:"nama"`
}
