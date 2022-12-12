package models

type Mahasiswa struct {
	Nim        string `json:"nim"`
	Nama       string `json:"nama"`
	Nomor_telp string `json:"nomor_telp"`
	Email      string `json:"email"`
	Angkatan   int    `json:"angkatan"`
	Role       string `json:"role"`
	Password   string `json:"password" gorm:"-"`
}

func (Mahasiswa) TableName() string {
	return "mahasiswa"
}
