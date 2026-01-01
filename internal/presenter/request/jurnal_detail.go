package request

import "github.com/google/uuid"

type JurnalDetailCreateRequest struct {
	IdJurnal uuid.UUID `validate:"required" json:"id_jurnal" gorm:"foreignKey:IdJurnal"`
	IdAkun   uuid.UUID `validate:"required" json:"id_akun"`
	Jenis    int       `validate:"required" json:"jenis"`
	Nominal  float64   `validate:"required" json:"nominal"`
}

type JurnalDetailUpdateRequest struct {
	Id       uuid.UUID
	IdJurnal uuid.UUID `validate:"required" json:"id_jurnal" gorm:"foreignKey:IdJurnal"`
	IdAkun   uuid.UUID `validate:"required" json:"id_akun"`
	Jenis    int       `validate:"required" json:"jenis"`
	Nominal  float64   `validate:"required" json:"nominal"`
}
