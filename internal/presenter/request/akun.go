package request

import "github.com/google/uuid"

type AkunCreateRequest struct {
	IdAkun        uuid.UUID `json:"id_akun"`
	IdAkunSubgrup uuid.UUID `validate:"required" json:"id_akun_subgrup"`
	Kode          string    `validate:"required,min=1,max=255" json:"kode"`
	Nama          string    `validate:"required,min=1,max=255" json:"nama"`
	Jenis         int       `validate:"required" json:"jenis"`
	SaldoAwal     float64   `json:"saldo_awal"`
	Status        int       `json:"status"`
}

type AkunUpdateRequest struct {
	Id            uuid.UUID
	IdAkun        uuid.UUID `json:"id_akun"`
	IdAkunSubgrup uuid.UUID `validate:"required" json:"id_akun_subgrup"`
	Kode          string    `validate:"required,min=1,max=255" json:"kode"`
	Nama          string    `validate:"required,min=1,max=255" json:"nama"`
	Jenis         int       `validate:"required" json:"jenis"`
	SaldoAwal     float64   `json:"saldo_awal"`
	Status        int       `json:"status"`
}
