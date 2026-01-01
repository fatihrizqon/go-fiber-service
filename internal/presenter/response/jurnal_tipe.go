package response

import "github.com/google/uuid"

type JurnalTipeResponse struct {
	Id   uuid.UUID `json:"id"`
	Nama string    `json:"nama"`
}
