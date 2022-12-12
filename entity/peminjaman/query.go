package peminjaman

var (
	GetDetailBukuDipinjamByIdPeminjaman = `
		SELECT p.id_buku, j.judul, j.tahun, j.penerbit FROM peminjaman_buku_map p 
		JOIN buku b on p.id_buku = b.id
		JOIN judul j on b.id_judul = j.id
		WHERE p.id_peminjaman = ? AND tanggal_kembali IS NULL
	`

	GetDetailBukuDipinjamByIdPeminjamanAll = `
		SELECT p.id_buku, j.judul, j.tahun, j.penerbit FROM peminjaman_buku_map p 
		JOIN buku b on p.id_buku = b.id
		JOIN judul j on b.id_judul = j.id
		WHERE p.id_peminjaman = ?
	`

	GetIdPeminjamanByNIM = `
		SELECT id, tanggal_peminjaman, status, tenggat_pengembalian FROM peminjaman
		WHERE nim_peminjaman = ? AND NOT status = 4
	`

	InsertPeminjaman = `
		INSERT INTO peminjaman (nim_peminjaman, tanggal_peminjaman, status)
		VALUES (?, ?, ?);
	`

	InsertPeminjamanIntoMap = `
		INSERT INTO peminjaman_buku_map (id_peminjaman, id_buku)
		VALUES (? , ?);
	`

	UpdateMapKembali = `
		UPDATE peminjaman_buku_map
		SET tanggal_kembali = ?
		WHERE id_buku = ? AND id_peminjaman = ?
	`

	CheckBukuAllKembali = `
		SELECT COUNT(*) FROM peminjaman_buku_map pbm
		WHERE pbm.id_peminjaman = ?
		AND tanggal_kembali IS NULL
	`

	UpdatePeminjamanKembali = `
		UPDATE peminjaman
		SET status = ?
		WHERE id = ?
	`

	GetAllPeminjaman = `
		SELECT * FROM peminjaman
	`

	GetDetailPeminjamanByID = `
		SELECT tanggal_peminjaman, status, nim_peminjaman, tenggat_pengembalian FROM peminjaman
		WHERE id = ?
	`
)
