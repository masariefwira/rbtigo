package buku

var (
	UpdateBukuDipinjam = `
		UPDATE buku SET status = ?
		WHERE id = ?
	`

	GetAllJudul = `
		SELECT j.id,
		 j.judul,
		 j.tahun,
		 j.penerbit,
		 COALESCE(j.penulis, ''),
		 COALESCE(j.filename, ''),
		 j.bahasa,
		 COALESCE(j.foto, ''),
		 COALESCE(j.jenis, 0),
		 COALESCE(j.id_kategori, 0)
		FROM judul j
	`

	GetBukuCountPerJudul = `
		SELECT COUNT(*)
		FROM buku
		WHERE id_judul = ?
	`

	UpdateJudulByJudulID = `
		UPDATE judul
		SET 
			judul = ?,
			tahun = ?,
			penerbit = ?,
			filename = ?,
			bahasa = ?,
			foto = ?,
			jenis = ?,
			id_kategori = ?
		WHERE id = ?
	`

	GetKategoriBuku = `
		SELECT kategori
		FROM kategori
		WHERE id = ?
	`

	GetAllKategori = `
		SELECT *
		FROM kategori
	`

	GetJudulByIDBuku = `
		SELECT id_judul from buku
		WHERE id = ?
	`

	GetJudulByID = `
		SELECT * from judul
		WHERE id = ?
	`

	GetAvailableJudulByIDJudul = `
		SELECT COUNT(*) FROM buku b
		WHERE id_judul = ? AND b.status = 1
	`

	AddBukuToCart = `
		INSERT INTO cart_item (id_judul, nim)
		VALUES (?,?)
	`

	DeleteItemFromCart = `
		DELETE FROM cart_item
		WHERE id_judul = ? AND nim = ?
	`

	GetCartItemsByNIM = `
		SELECT id_judul FROM cart_item ci
		WHERE nim = ?
	`

	GetAvailableIDBukuByIDJudul = `
		SELECT b.id FROM buku b 
		WHERE id_judul = ? AND b.status = 1
	`

	GetDetailBukuByJudulID = `
		SELECT id, status from buku
		WHERE id_judul = ?
	`

	InsertBukuIndividual = `
		INSERT INTO buku (id, status, id_judul)
		VALUES (?, ?, ?)
	`

	DeleteBuku = `
		DELETE FROM judul
		WHERE id = ?
	`
)