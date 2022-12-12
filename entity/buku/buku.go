package buku

import (
	"errors"
	"fmt"

	"github.com/ikalkali/rbti-go/entity/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BukuEntityInterface interface {
	UpdateBukuDipinjam(idBuku int, status int, tx *gorm.DB) error
	GetAllJudul(filter models.FilterBuku) ([]models.Judul, error)
	GetBukuCountPerJudul(idJudul int) (int, error)
	InsertJudulBaru(input models.Judul, tx *gorm.DB) (int, error)
	BatchInsertBukuIndividualBaru(input []models.Buku, tx *gorm.DB) error
	UpdateJudul(input models.Judul, tx *gorm.DB) error
	GetJudulByID(idJudul int) (models.Judul, error)
	GetKategoriBuku(idKategori int) (string, error)
	GetAllKategori() ([]models.Kategori, error)
	GetIDJudulByBukuID(idBuku int) (int, error)
	GetAvailableJudulByIDJudul(idJudul int) (int, error)
	AddBukuToCart(idBuku int, nim string, tx *gorm.DB) error
	DeleteItemFromCart(idBuku int, nim string, tx *gorm.DB) error
	GetCartItemsByNIM(nim string) ([]int, error)
	GetDetailBukuByJudulID(idJudul int) ([]models.DetailBukuPerJudul, error)
	InsertIDBuku(idBuku int, idJudul int, tx *gorm.DB) error
	DeleteBuku(id int, tx *gorm.DB) error
	GetAvailableIDBukuByIDJudul(idJudul int) (int, error)
}

type buku struct {
	db *gorm.DB
}

func NewEntity(db *gorm.DB) *buku {
	return &buku{db}
}

func (b *buku) UpdateBukuDipinjam(idBuku int, status int, tx *gorm.DB) error {
	if err := tx.Exec(UpdateBukuDipinjam, status, idBuku).Error; err != nil {
		return err
	}

	return nil
}

func (b *buku) GetAllJudul(filter models.FilterBuku) ([]models.Judul, error) {
	var (
		resp      []models.Judul
		limit     = filter.Limit
		offset    = filter.Offset
		tempQuery string
		where     string
		count     int64
	)

	limitQuery := "LIMIT ? OFFSET ?"

	if filter.IDJudul != 0 {
		tempQuery = fmt.Sprintf("WHERE id = %v", filter.IDJudul)
	} else {
		tempQuery = getAllJudulQueryBuilder(filter)
	}

	c := "SELECT COUNT(*) from judul"

	finalQuery := fmt.Sprintf("%v %v %v", GetAllJudul, tempQuery, limitQuery)

	fmt.Println(finalQuery)
	countQuery := fmt.Sprintf("%v %v %v", c, where, tempQuery)

	rows, err := b.db.Raw(finalQuery, limit, offset).Rows()
	if err != nil {
		return nil, err
	}

	err = b.db.Raw(countQuery).Scan(&count).Error
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := models.Judul{}
		err = rows.Scan(&temp.Id,
			&temp.Judul,
			&temp.Tahun,
			&temp.Penerbit,
			&temp.Penulis,
			// &temp.Filename,
			&temp.Bahasa,
			// &temp.Foto,
			&temp.Jenis,
			&temp.IDKategori,
		)

		if err != nil {
			logrus.Error(err.Error())
		}

		temp.Count = count

		resp = append(resp, temp)
	}

	return resp, nil
}

func (b *buku) GetBukuCountPerJudul(idJudul int) (int, error) {
	var (
		resp int
	)

	err := b.db.Raw(GetBukuCountPerJudul, idJudul).Scan(&resp).Error
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (b *buku) GetKategoriBuku(idKategori int) (string, error) {
	var (
		resp string
	)

	err := b.db.Raw(GetKategoriBuku, idKategori).Scan(&resp).Error
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (b *buku) InsertJudulBaru(input models.Judul, tx *gorm.DB) (int, error) {
	if err := tx.Create(&input).Error; err != nil {
		return 0, err
	}

	return input.Id, nil
}

func (b *buku) BatchInsertBukuIndividualBaru(input []models.Buku, tx *gorm.DB) error {
	if err := tx.Create(&input).Error; err != nil {
		return err
	}

	return nil
}

func (b *buku) UpdateJudul(input models.Judul, tx *gorm.DB) error {
	var (
		judul    = input.Judul
		tahun    = input.Tahun
		penerbit = input.Penerbit
		// filename    = input.Filename
		bahasa = input.Bahasa
		// foto        = input.Foto
		jenis       = input.Jenis
		id_kategori = input.IDKategori
		id          = input.Id
	)

	if err := tx.Exec(UpdateJudulByJudulID,
		judul,
		tahun,
		penerbit,
		// filename,
		bahasa,
		// foto,
		jenis,
		id_kategori,
		id,
	).Error; err != nil {
		return err
	}

	return nil
}

func (b *buku) GetJudulByID(idJudul int) (models.Judul, error) {
	var (
		resp models.Judul
	)

	err := b.db.Where("id = ?", idJudul).Find(&resp).Error
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (b *buku) GetAllKategori() ([]models.Kategori, error) {
	var (
		resp []models.Kategori
	)

	rows, err := b.db.Raw(GetAllKategori).Rows()
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		temp := models.Kategori{}
		rows.Scan(&temp.Id, &temp.Kategori)
		resp = append(resp, temp)
	}

	return resp, nil
}

func (b *buku) GetIDJudulByBukuID(idBuku int) (int, error) {
	var (
		resp int
	)

	err := b.db.Raw(GetJudulByIDBuku, idBuku).Scan(&resp).Error
	if err != nil {
		return resp, err
	}

	if resp == 0 {
		return resp, errors.New("no judul with the given ID found, please provide a correct ID")
	}

	return resp, nil
}

// return count available
func (b *buku) GetAvailableJudulByIDJudul(idJudul int) (int, error) {
	var (
		resp int
	)

	err := b.db.Raw(GetAvailableJudulByIDJudul, idJudul).Scan(&resp).Error
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (b *buku) AddBukuToCart(idBuku int, nim string, tx *gorm.DB) error {
	err := tx.Exec(AddBukuToCart, idBuku, nim).Error
	if err != nil {
		return err
	}

	return nil
}

func (b *buku) DeleteItemFromCart(idBuku int, nim string, tx *gorm.DB) error {
	err := tx.Exec(DeleteItemFromCart, idBuku, nim).Error
	if err != nil {
		return err
	}

	return nil
}

func (b *buku) GetCartItemsByNIM(nim string) ([]int, error) {
	var (
		resp []int
	)

	err := b.db.Raw(GetCartItemsByNIM, nim).Scan(&resp).Error
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// return ID Buku
func (b *buku) GetAvailableIDBukuByIDJudul(idJudul int) (int, error) {
	var (
		resp int
	)

	err := b.db.Raw(GetAvailableIDBukuByIDJudul, idJudul).Scan(&resp).Error
	if err != nil {
		return resp, err
	}

	return resp, nil

}

func (b *buku) GetDetailBukuByJudulID(idJudul int) ([]models.DetailBukuPerJudul, error) {
	var (
		resp []models.DetailBukuPerJudul
	)

	rows, err := b.db.Raw(GetDetailBukuByJudulID, idJudul).Rows()
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var temp models.DetailBukuPerJudul
		rows.Scan(&temp.IDBuku, &temp.Status)
		resp = append(resp, temp)
	}

	return resp, nil
}

func (b *buku) InsertIDBuku(idBuku int, idJudul int, tx *gorm.DB) error {
	err := tx.Exec(InsertBukuIndividual, idBuku, 1, idJudul).Error
	if err != nil {
		return err
	}

	return nil
}

func (b *buku) DeleteBuku(id int, tx *gorm.DB) error {
	err := tx.Exec(DeleteBuku, id).Error
	if err != nil {
		return err
	}

	return nil
}
