package response

import (
	"time"

	"github.com/google/uuid"
)

type JurnalDetailResponse struct {
	Id        uuid.UUID `json:"id"`
	IdJurnal  uuid.UUID `json:"id_jurnal"`
	IdAkun    uuid.UUID `json:"id_akun"`
	NamaAkun  string    `json:"nama_akun"`
	Jenis     int       `json:"jenis"`
	Nominal   float64   `json:"nominal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
