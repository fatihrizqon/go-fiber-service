package request

import (
	"time"

	"github.com/google/uuid"
)

type JurnalCreateRequest struct {
	IdJurnalTipe uuid.UUID                   `validate:"required" json:"id_jurnal_tipe"`
	Tanggal      time.Time                   `json:"tanggal"`
	Nota         string                      `json:"nota"`
	Keterangan   string                      `validate:"required" json:"keterangan"`
	CreatedBy    uuid.UUID                   `json:"created_by"`
	JurnalDetail []JurnalDetailCreateRequest `validate:"required" json:"jurnal_detail"`
}

type JurnalUpdateRequest struct {
	Id           uuid.UUID
	IdJurnalTipe uuid.UUID                   `validate:"required" json:"id_jurnal_tipe"`
	Tanggal      time.Time                   `json:"tanggal"`
	Nota         string                      `json:"nota"`
	Keterangan   string                      `validate:"required" json:"keterangan"`
	UpdatedBy    uuid.UUID                   `json:"updated_by"`
	JurnalDetail []JurnalDetailUpdateRequest `validate:"required" json:"jurnal_detail"`
}
