package mahasiswa

var (
	InsertMahasiswaBaru = `
		INSERT INTO mahasiswa (nim,nama,nomor_telp,email,CAST (angkatan AS INTEGER))
		VALUES (?,?,?,?,?)
	`

	InsertMahasiswaBaruSignup = `
		INSERT INTO mahasiswa (nim,nama,nomor_telp,email,angkatan, password, role)
		VALUES (?,?,?,?,?,?,?)
	`

	UpdatePasswordMahasiswa = `
		UPDATE mahasiswa
		SET
			email = ?,
			password = ?
		WHERE nim = ?
	`

	CheckMahasiswaExists = `
		SELECT COUNT(*) FROM mahasiswa
		WHERE nim = ?
	`

	CheckMahasiswaExistEligible = `
		SELECT COUNT(*) FROM mahasiswa
		WHERE nim = ? AND password is not null
	`

	GetPasswordByEmail = `
		SELECT password FROM mahasiswa
		WHERE email = ?
	`

	GetDetailMahasiswaByEmail = `
		SELECT nim, nama, nomor_telp, email, angkatan, role
		FROM mahasiswa
		WHERE email = ?
	`
)
