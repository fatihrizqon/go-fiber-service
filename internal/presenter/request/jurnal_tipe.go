package request

import "github.com/google/uuid"

type JurnalTipeCreateRequest struct {
	Nama string `validate:"required,min=1,max=255" json:"nama"`
}

type JurnalTipeUpdateRequest struct {
	Id   uuid.UUID
	Nama string `validate:"required,min=1,max=255" json:"nama"`
}
