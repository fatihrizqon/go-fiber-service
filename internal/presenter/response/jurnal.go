package response

import (
	"time"

	"github.com/google/uuid"
)

type JurnalResponse struct {
	Id           uuid.UUID              `json:"id"`
	IdJurnalTipe uuid.UUID              `json:"id_jurnal_tipe"`
	Jenis        int                    `json:"jenis"`
	Tanggal      time.Time              `json:"tanggal"`
	Nota         string                 `json:"nota"`
	Keterangan   string                 `json:"keterangan"`
	CreatedBy    uuid.UUID              `json:"created_by"`
	UpdatedBy    uuid.UUID              `json:"updated_by"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	JurnalDetail []JurnalDetailResponse `json:"jurnal_detail"`
}

type ReportJurnalResponse struct {
	IdJurnalTipe uuid.UUID `json:"id_jurnal_tipe"`
	Kategori     string    `json:"kategori"`
	IdAkun       uuid.UUID `json:"id_akun"`
	Kode         string    `json:"kode"`
	Akun         string    `json:"akun"`
	Nota         string    `json:"nota"`
	Jenis        int       `json:"jenis"`
	Nominal      float64   `json:"nominal"`
	Tanggal      time.Time `json:"tanggal"`
	Keterangan   string    `json:"keterangan"`
}
