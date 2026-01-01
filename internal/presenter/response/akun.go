package response

import (
	"time"

	"github.com/google/uuid"
)

type AkunResponse struct {
	Id              uuid.UUID `json:"id"`
	IdAkunSubgrup   uuid.UUID `json:"id_akun_subgrup"`
	NamaAkunSubgrup string    `json:"nama_akun_subgrup"`
	Kode            string    `json:"kode"`
	Nama            string    `json:"nama"`
	Jenis           int       `json:"jenis"`
	SaldoAwal       float64   `json:"saldo_awal"`
	Status          int       `json:"status"`
}

type BukuBesarResponse struct {
	IdAkun      uuid.UUID `json:"id_akun"`
	KodeGrup    string    `json:"kode_grup"`
	Grup        string    `json:"grup"`
	KodeSubgrup string    `json:"kode_subgrup"`
	Subgrup     string    `json:"subgrup"`
	Kode        string    `json:"kode"`
	Akun        string    `json:"akun"`
	SaldoAwal   float64   `json:"saldo_awal"`
	TotalDebit  float64   `json:"total_debit"`
	TotalKredit float64   `json:"total_kredit"`
	SaldoAkhir  float64   `json:"saldo_akhir"`
}

type BukuBesarDetailResponse struct {
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

type TransaksiAkunResponse struct {
	IdAkun      uuid.UUID `json:"id_akun"`
	KodeGrup    string    `json:"kode_grup"`
	Grup        string    `json:"grup"`
	KodeSubgrup string    `json:"kode_subgrup"`
	Subgrup     string    `json:"subgrup"`
	Kode        string    `json:"kode"`
	Akun        string    `json:"akun"`
	Jenis       int       `json:"jenis"`
	SaldoAwal   float64   `json:"saldo_awal"`
	TotalDebit  float64   `json:"total_debit"`
	TotalKredit float64   `json:"total_kredit"`
	SaldoAkhir  float64   `json:"saldo_akhir"`
	Bulan       string    `json:"bulan"`
	Tahun       string    `json:"tahun"`
}

type StatistikLabaRugiResponse struct {
	Bulan      string  `json:"bulan"`
	Tahun      string  `json:"tahun"`
	Pendapatan float64 `json:"pendapatan"`
	Biaya      float64 `json:"biaya"`
	LabaRugi   float64 `json:"laba_rugi"`
	Akumulatif float64 `json:"akumulatif"`
}
