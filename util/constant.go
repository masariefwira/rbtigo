package util

import "time"

// Postgres DSN
var (
	POSTGRES_DSN string = "user=postgres password=root dbname=rbtigo port=5432 sslmode=disable"
)

// Status
var (
	STATUS_TERSEDIA               = 1
	STATUS_DIPINJAM               = 2
	STATUS_DIPINJAM_BELUM_DIAMBIL = 3

	STATUS_KEMBALI         = 4
	STATUS_KEMBALI_PARTIAL = 5
	STATUS_CANCEL          = 6
)

var (
	BISA_DIPINJAM       = 1
	TIDAK_BISA_DIPINJAM = 2
)

// Peminjaman
var (
	DURASI_PEMINJAMAN = (24 * time.Hour) * 14 // 14 hari
	DENDA_PER_HARI    = 5000                  // Rp. 5000
	MASIH_DIPINJAM    = "dipinjam"
	SEMUA_BUKU        = "semua"
)
