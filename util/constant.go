package util

import "time"

// Postgres DSN
var (
	POSTGRES_DSN string = "user=postgres password=root dbname=rbtigo port=5432 sslmode=disable"
	// POSTGRES_DSN string = "postgres://uzhvdklucnaowc:19c420ecf6ef5002e5b9adaf895e262716b2bb0a5f307de28dbe7feb4258a764@ec2-34-194-40-194.compute-1.amazonaws.com:5432/d7hgf4q0rjb1tr"
	// POSTGRES_DSN_STAGING string = "host=103.193.14.37 user=postgres password=ikalkali dbname=rbti port=5432 sslmode=disable"
	// POSTGRES_DSN_STAGING string = "postgres://uzhvdklucnaowc:19c420ecf6ef5002e5b9adaf895e262716b2bb0a5f307de28dbe7feb4258a764@ec2-34-194-40-194.compute-1.amazonaws.com:5432/d7hgf4q0rjb1tr"
	// POSTGRES_DSN_STAGING string = "host=ec2-34-194-40-194.compute-1.amazonaws.com user=uzhvdklucnaowc password=19c420ecf6ef5002e5b9adaf895e262716b2bb0a5f307de28dbe7feb4258a764 dbname=d7hgf4q0rjb1tr port=5432 sslmode=disable"
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
