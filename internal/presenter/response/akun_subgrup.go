package response

import "github.com/google/uuid"

type AkunSubgrupResponse struct {
	Id           uuid.UUID `json:"id"`
	IdAkunGrup   uuid.UUID `json:"id_akun_grup"`
	NamaAkunGrup string    `json:"nama_akun_grup"`
	Kode         string    `json:"kode"`
	Nama         string    `json:"nama"`
}
