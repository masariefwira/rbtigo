package laporan

import (
	"fmt"

	"github.com/ikalkali/rbti-go/entity/models"
	"gorm.io/gorm"
)

type LaporanEntityInterface interface {
	GetAllLaporan(input models.LaporanFilter) ([]models.LaporanDB, error)
	GetLaporanCount() (int, error)
	InsertLaporan(input models.LaporanDB, tx *gorm.DB) (int, error)
	DeleteLaporan(id int, tx *gorm.DB) error
	InsertPaper(input models.Paper, tx *gorm.DB) (int, error)
	SearchPaper(query string, jenis int) ([]models.Paper, int, error)
	InsertArtikelOrMakalah(input models.LaporanDB, tx *gorm.DB) (int, error)
	InsertKaryaTulis(input models.KaryaTulis, tx *gorm.DB) (int, error)
	GetIDKaryaTulisByIDPaper(idPaper int) ([]int, error)
	GetAllPaper(filter models.PaperFilter) ([]models.Paper, int, error)
	GetAllArtikel(filter models.PaperFilter) ([]models.Artikel, int, error)
	GetDetailPaperByID(idPaper int) (models.Paper, error)
}

type laporan struct {
	db *gorm.DB
}

func NewEntity(db *gorm.DB) *laporan {
	return &laporan{db}
}

func (l *laporan) GetAllLaporan(input models.LaporanFilter) ([]models.LaporanDB, error) {
	var (
		resp []models.LaporanDB
	)
	rows, err := l.db.Raw(GetAllLaporan, input.Limit, input.Offset).Rows()
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		temp := models.LaporanDB{}
		err = rows.Scan(
			&temp.Id,
			// &temp.Tahun,
			&temp.NIM,
			&temp.IdKategori,
			&temp.Jenis,
			&temp.Judul,
		)

		if err != nil {
			return resp, err
		}

		resp = append(resp, temp)
	}

	return resp, nil
}

func (l *laporan) GetLaporanCount() (int, error) {
	var (
		resp int
	)

	err := l.db.Raw(GetLaporanCount).Scan(&resp).Error
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (l *laporan) InsertLaporan(input models.LaporanDB, tx *gorm.DB) (int, error) {
	var id int
	err := tx.Raw(InsertLaporan,
		// input.Tahun,
		input.NIM,
		input.IdKategori,
		input.Jenis,
		input.Judul,
	).Scan(&id).Error

	if err != nil {
		return id, err
	}

	return id, nil
}

func (l *laporan) InsertArtikelOrMakalah(input models.LaporanDB, tx *gorm.DB) (int, error) {
	var id int
	err := tx.Raw(InsertMakalahOrArtikel,
		// input.Tahun,
		input.NIM,
		input.IdKategori,
		input.Jenis,
		input.Judul,
		input.IdPaper,
	).Scan(&id).Error

	if err != nil {
		return id, err
	}

	return id, nil
}

func (l *laporan) DeleteLaporan(id int, tx *gorm.DB) error {
	err := tx.Exec(DeleteLaporan, id).Error
	if err != nil {
		return err
	}

	return nil
}

func (l *laporan) InsertPaper(input models.Paper, tx *gorm.DB) (int, error) {
	var id int
	err := tx.Raw(InsertPaper, input.Judul, input.Volume, input.Jenis, input.Tahun).Scan(&id).Error
	if err != nil {
		return id, err
	}

	return id, nil
}

func (l *laporan) SearchPaper(query string, jenis int) ([]models.Paper, int, error) {
	var resp []models.Paper
	var count int
	fmt.Printf(SearchPaperByJudul, query, jenis)
	rows, err := l.db.Raw(fmt.Sprintf(SearchPaperByJudul, query, jenis)).Rows()
	if err != nil {
		return resp, count, err
	}
	for rows.Next() {
		temp := models.Paper{}
		err = rows.Scan(&temp.Id, &temp.Judul,
			&temp.Volume,
			&temp.Jenis,
			&temp.Tahun)
		if err != nil {
			return resp, count, err
		}

		resp = append(resp, temp)
		count = count + 1
	}
	if err != nil {
		return resp, count, err
	}

	return resp, count, nil
}

func (l *laporan) GetAllPaper(filter models.PaperFilter) ([]models.Paper, int, error) {
	var resp []models.Paper
	var count int
	rows, err := l.db.Raw(GetAllPaper, filter.Jenis, filter.Limit, filter.Offset).Rows()
	if err != nil {
		return resp, count, err
	}

	for rows.Next() {
		temp := models.Paper{}
		err = rows.Scan(&temp.Id, &temp.Judul,
			&temp.Volume,
			&temp.Jenis,
			&temp.Tahun)
		if err != nil {
			return resp, count, err
		}

		resp = append(resp, temp)
		count = count + 1
	}
	if err != nil {
		return resp, count, err
	}

	return resp, count, nil
}

func (l *laporan) GetAllArtikel(filter models.PaperFilter) ([]models.Artikel, int, error) {
	var resp []models.Artikel
	var count int
	rows, err := l.db.Raw(GetAllArtikel, filter.Limit, filter.Offset).Rows()
	if err != nil {
		return resp, count, err
	}

	for rows.Next() {
		temp := models.Artikel{}
		err = rows.Scan(
			&temp.Id,
			&temp.IdPaper,
			&temp.IdKategori,
			&temp.Judul,
			&temp.Tahun,
			&temp.Penulis,
		)
		if err != nil {
			return resp, count, err
		}

		resp = append(resp, temp)
		count = count + 1
	}
	if err != nil {
		return resp, count, err
	}

	return resp, count, nil
}

func (l *laporan) GetDetailPaperByID(idPaper int) (models.Paper, error) {
	var resp models.Paper
	err := l.db.Raw(GetDetailPaperByID, idPaper).Scan(&resp).Error
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (l *laporan) InsertKaryaTulis(input models.KaryaTulis, tx *gorm.DB) (int, error) {
	var id int
	err := tx.Raw(InsertKaryaTulis,
		input.IDPaper,
		input.IDKategori,
		input.Judul,
		input.Tahun,
		input.Penulis,
	).Scan(&id).Error

	if err != nil {
		return id, err
	}

	return id, nil
}

func (l *laporan) GetIDKaryaTulisByIDPaper(idPaper int) ([]int, error) {
	var resp []int
	rows, err := l.db.Raw(GetIDKaryaTulisByIDPaper, idPaper).Rows()
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var temp int
		err = rows.Scan(&temp)
		if err != nil {
			return resp, err
		}

		resp = append(resp, temp)
	}

	return resp, nil
}
