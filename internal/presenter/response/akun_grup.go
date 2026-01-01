package response

import "github.com/google/uuid"

type AkunGrupResponse struct {
	Id   uuid.UUID `json:"id"`
	Kode string    `json:"kode"`
	Nama string    `json:"nama"`
}
