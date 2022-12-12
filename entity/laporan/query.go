package laporan

var (
	GetAllLaporan = `
		SELECT * FROM laporan
		LIMIT ? OFFSET ?
	`

	GetLaporanCount = `
		SELECT COUNT(*) FROM laporan
	`
	InsertLaporan = `
		INSERT INTO laporan (tahun, nim_penulis, id_kategori, jenis, judul)
		VALUES (?, ?, ?, ?, ?) RETURNING id
	`

	InsertMakalahOrArtikel = `
		INSERT INTO laporan (tahun, nim_penulis, id_kategori, jenis, judul, id_paper)
		VALUES (?, ?, ?, ?, ?, ?) RETURNING id
	`

	DeleteLaporan = `
		DELETE FROM laporan
		where id = ?
	`

	InsertPaper = `
		INSERT INTO paper (judul, volume, jenis, tanggal_rilis)
		VALUES (?, ?, ?, ?) RETURNING id
	`

	SearchPaperByJudul = `
		SELECT id, judul, volume, jenis, tanggal_rilis FROM paper p 
		WHERE LOWER(judul) LIKE LOWER('%s%%') AND jenis = %v
	`

	InsertKaryaTulis = `
		INSERT INTO karya_tulis (id_paper, id_kategori, judul, tahun, penulis)
		VALUES (?, ?, ?, ?, ?) RETURNING id
	`

	GetIDKaryaTulisByIDPaper = `
		SELECT id FROM karya_tulis
		WHERE id_paper = ?
	`

	GetAllPaper = `
		SELECT id, judul, volume, jenis, tanggal_rilis
		FROM paper p
		WHERE jenis = ?
		LIMIT ? OFFSET ?
	`

	GetAllArtikel = `
		SELECT id, id_paper, id_kategori, judul, tahun, penulis
		FROM karya_tulis kt
		LIMIT ? OFFSET ?
	`

	GetDetailPaperByID = `
		SELECT id, judul, volume, jenis, tanggal_rilis
		FROM paper p
		WHERE id = ?
	`
)
