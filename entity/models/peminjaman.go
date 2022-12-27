package models

import "time"

type DetailPeminjaman struct {
	Judul    string `json:"judul"`
	Tahun    int    `json:"tahun"`
	Penerbit string `json:"penerbit"`
	IDBuku   int    `json:"id_buku"`
}

type DetailPeminjamanRequest struct {
	IDPeminjaman int `json:"id_peminjaman"`
}

type Peminjaman struct {
	IDPeminjaman        int                `json:"id_peminjaman" gorm:"-"`
	TanggalPeminjaman   time.Time          `json:"tanggal_peminjaman" db:"tanggal_peminjaman"`
	TenggatPengembalian time.Time          `json:"tenggat_pengembalian" db:"tenggat_pengembalian"`
	BukuDipinjam        []DetailPeminjaman `json:"buku_dipinjam" gorm:"-"`
	NIM                 string             `json:"nim" gorm:"nim_peminjaman" db:"nim_peminjaman"`
	Denda               int64              `json:"denda" gorm:"-"`
	DetailDenda         string             `json:"detail_denda" gorm:"-"`
	Status              int                `json:"status" db:"status"`
}

type DetailPeminjamanByNIM struct {
	Peminjaman []Peminjaman `json:"peminjaman"`
	TotalDenda int64        `json:"total_denda"`
}

type DetailPeminjamanByNIMRequest struct {
	NIM string `json:"nim"`
}

type InputPeminjamanRequest struct {
	ID      int    `json:"id"`
	NIM     string `json:"nim"`
	IDBuku  []int  `json:"id_buku"`
	IDJudul []int  `json:"id_judul"`
	Source  string `json:"source"`
}

type InputPengembalianRequest struct {
	IDPeminjaman int   `json:"id_peminjaman"`
	IDBuku       []int `json:"id_buku"`
}

type InputPengembalianData struct {
	Data []InputPengembalianRequest `json:"data"`
}

type PeminjamanDB struct {
	ID                  int       `db:"id"`
	TanggalPeminjaman   time.Time `db:"tanggal_peminjaman"`
	Status              int       `db:"status"`
	NIM                 string    `db:"nim_peminjaman" gorm:"column:nim_peminjaman"`
	TenggatPengembalian time.Time `db:"tenggat_pengembalian"`
}

type PeminjamanMap struct {
	ID             int        `db:"id"`
	TanggalKembali *time.Time `db:"tanggal_kembali"`
	IDBuku         int        `db:"id_buku"`
	IDPeminjaman   int        `db:"id_peminjaman"`
}

type PeminjamanFilter struct {
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
	Telat  string                `json:"telat"`
	Status []int                 `json:"status"`
	Waktu  PeminjamanFilterWaktu `json:"waktu"`
}

type PeminjamanFilterWaktu struct {
	Waktu    string `json:"waktu"`
	Operator string `json:"operator"`
}

type PeminjamanResponse struct {
	ID                  int       `json:"id" gorm:"column:id"`
	TanggalPeminjaman   time.Time `json:"tanggal_peminjaman" gorm:"column:tanggal_peminjaman"`
	Status              int       `json:"status" gorm:"column:status"`
	NIM                 string    `json:"nim_peminjaman" gorm:"column:nim_peminjaman"`
	Nama                string    `json:"nama"`
	TenggatPengembalian time.Time `json:"tenggat_pengembalian" gorm:"column:tenggat_pengembalian"`
	Keterlambatan       string    `json:"keterlambatan"`
	Denda               int64     `json:"denda"`
	Count               int64     `json:"count"`

	// for notify
	DetailPeminjaman []DetailPeminjaman `json:"detail_peminjaman,omitempty"`
	Email            string             `json:"email,omitempty"`
}

type DetailBukuPeminjamanFilter struct {
	IDPeminjaman int
	Ketersediaan string
}

func (PeminjamanDB) TableName() string {
	return "peminjaman"
}

func (PeminjamanMap) TableName() string {
	return "peminjaman_buku_map"
}
