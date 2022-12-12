package models

type Judul struct {
	Id       int    `json:"id"`
	Judul    string `json:"judul"`
	Tahun    int    `json:"tahun"`
	Penulis  string `json:"penulis"`
	Penerbit string `json:"penerbit"`
	// Filename    *string              `json:"filename"`
	Bahasa string `json:"bahasa"`
	// Foto        *string              `json:"foto"`
	Jenis       int                  `json:"jenis"`
	IDKategori  int                  `json:"id_kategori"`
	Count       int64                `json:"count" gorm:"-"`
	DetailBuku  []DetailBukuPerJudul `json:"detail_buku" gorm:"-"`
	IDBuku      []int                `json:"id_buku,omitempty" gorm:"-"`
	IsAvailable bool                 `json:"is_available" gorm:"-"`
}

type DetailBukuPerJudul struct {
	IDBuku int `json:"id_buku"`
	Status int `json:"status"`
}

type JudulElastic struct {
	Id         int         `json:"id"`
	Judul      string      `json:"judul"`
	Tahun      int         `json:"tahun"`
	Penulis    string      `json:"penulis"`
	Tipe       string      `json:"tipe"`
	Kategori   interface{} `json:"kategori"`
	IDKategori int         `json:"id_kategori"`

	// buku fields
	Penerbit       string      `json:"penerbit,omitempty"`
	Bahasa         string      `json:"bahasa,omitempty"`
	Jenis          interface{} `json:"jenis,omitempty"`
	JumlahTotal    int         `json:"jumlah_total,omitempty"`
	JumlahTersedia int         `json:"jumlah_tersedia,omitempty"`

	// laporan fields
	NIM string `json:"nim,omitempty"`

	// jurnal fields
	IdPaper int `json:"id_paper,omitempty"`
}

type SearchBukuRequest struct {
	Query  string `json:"query"`
	Filter string `json:"filter"`
	From   int    `json:"from"`
	Size   int    `json:"size"`
}

type SearchBukuResponse struct {
	Data   []JudulElastic `json:"data"`
	Errors []string       `json:"errors"`
	Count  int            `json:"count"`
}

type Buku struct {
	Id      int
	Status  int
	IdJudul int
}

type GetAllJudulResponseDetail struct {
	Judul  Judul `json:"judul"`
	Jumlah int   `json:"jumlah"`
}

type GetAllJudulResponse struct {
	Data   []GetAllJudulResponseDetail `json:"data"`
	Errors []string                    `json:"errors"`
}

type InputJudulBaruRequest struct {
	DetailJudul Judul `json:"judul"`
	Jumlah      int   `json:"jumlah"`
}

type UpdateJudulBukuRequest struct {
	Id          int   `json:"id"`
	DetailJudul Judul `json:"judul"`
}

type Kategori struct {
	Id       int    `json:"id"`
	Kategori string `json:"kategori"`
}

type FilterBuku struct {
	Kategori []int  `json:"kategori"`
	Jenis    []int  `json:"jenis"`
	Judul    string `json:"judul"`
	Penerbit string `json:"penerbit"`
	Penulis  string `json:"penulis"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	IDJudul  int    `json:"id_judul"`
}

type CartRequest struct {
	IdJudul []int  `json:"id_judul"`
	NIM     string `json:"nim"`
	Action  string `json:"action"`
}

type CartResponse struct {
	Judul []Judul `json:"judul"`
}

type ElasticFilter struct {
	ID           int    `json:"id"`
	Query        string `json:"query"`
	From         int    `json:"from"`
	Size         int    `json:"size"`
	Jenis        string `json:"jenis"`
	Kategori     string `json:"kategori"`
	IDKategori   []int  `json:"id_kategori"`
	JenisPinjam  int    `json:"jenis_pinjam"`
	BisaDipinjam bool   `json:"bisa_dipinjam"`
}

func (Judul) TableName() string {
	return "judul"
}

func (Buku) TableName() string {
	return "buku"
}
