package models

import "time"

type LaporanFilter struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type LaporanDB struct {
	Id         int    `json:"id"`
	Tahun      int    `json:"tahun"`
	NIM        string `json:"nim"`
	IdKategori int    `json:"id_kategori"`
	Jenis      string `json:"jenis"`
	Judul      string `json:"judul"`
	IdPaper    int    `json:"id_paper"`
}

type Artikel struct {
	Id         int    `json:"id"`
	Tahun      int    `json:"tahun"`
	Penulis    string `json:"penulis"`
	IdKategori int    `json:"id_kategori"`
	Jenis      string `json:"jenis"`
	Judul      string `json:"judul"`
	IdPaper    int    `json:"id_paper"`
	JudulInduk string `json:"judul_induk"`
}

type LaporanES struct {
	Id         int    `json:"id"`
	Tahun      int    `json:"tahun"`
	NIM        string `json:"nim"`
	Penulis    string `json:"penulis"`
	IdKategori int    `json:"id_kategori"`
	Kategori   string `json:"kategori"`
	Jenis      string `json:"jenis"`
	Judul      string `json:"judul"`
}

type Paper struct {
	Id     int       `json:"id"`
	Volume string    `json:"volume"`
	Judul  string    `json:"judul"`
	Jenis  int       `json:"jenis"`
	Tahun  time.Time `json:"tahun"`
}

type KaryaTulis struct {
	Id      int    `json:"id"`
	Judul   string `json:"judul"`
	Penulis string `json:"penulis"`
	Tahun   int    `json:"tahun"`

	// foreign keys
	IDPaper    int `json:"id_paper"`
	IDKategori int `json:"id_kategori"`
}

type PaperFilter struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Jenis  int `json:"jenis"`
}
