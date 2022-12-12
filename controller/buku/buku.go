package buku

import (
	"fmt"
	"strings"
	"time"

	esBuku "github.com/ikalkali/rbti-go/elastic/buku"
	esLaporan "github.com/ikalkali/rbti-go/elastic/laporan"
	"github.com/ikalkali/rbti-go/entity/buku"
	"github.com/ikalkali/rbti-go/entity/laporan"
	"github.com/ikalkali/rbti-go/entity/models"
	"github.com/ikalkali/rbti-go/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BukuControllerInterface interface {
	GetAllJudul(filter models.FilterBuku) ([]models.GetAllJudulResponseDetail, error)
	InputJudulBaru(input models.InputJudulBaruRequest) error
	UpdateJudulBuku(input models.Judul) error
	SyncBukuDBWithES() error
	SearchBuku(input models.ElasticFilter) ([]models.JudulElastic, int, error)
	GetAllKategori() ([]models.Kategori, error)
	GetJudulByBukuID(idBuku int) (models.Judul, error)
	GetCartItemsByNIM(nim string) (models.CartResponse, error)
	GetJudulByIDElastic(input models.ElasticFilter) (models.JudulElastic, error)
	DeleteBukuOrLaporan(input models.ElasticFilter) (error)
}

type bukuController struct {
	bukuDb buku.BukuEntityInterface
	bukuEs esBuku.BukuElasticInterface
	laporanEs esLaporan.LaporanElasticInterface
	laporanDB laporan.LaporanEntityInterface
	db     *gorm.DB
}

func NewController(
	bukuDb buku.BukuEntityInterface,
	db *gorm.DB,
	bukuEs esBuku.BukuElasticInterface,
	laporanEs esLaporan.LaporanElasticInterface,
	laporanDB laporan.LaporanEntityInterface,
) *bukuController {
	return &bukuController{
		bukuDb: bukuDb,
		db:     db,
		bukuEs: bukuEs,
		laporanEs: laporanEs,
		laporanDB: laporanDB,
	}
}

func (b *bukuController) GetAllJudul(filter models.FilterBuku) ([]models.GetAllJudulResponseDetail, error) {
	var (
		resp []models.GetAllJudulResponseDetail
	)

	// Get All Judul
	judul, err := b.bukuDb.GetAllJudul(filter)
	if err != nil {
		return nil, err
	}

	// Populate buku with stock
	for _, book := range judul {
		temp := models.GetAllJudulResponseDetail{}
		idJudul := book.Id
		count, err := b.bukuDb.GetBukuCountPerJudul(idJudul)

		if err != nil {
			return nil, err
		}

		detailBuku, err := b.bukuDb.GetDetailBukuByJudulID(idJudul)
		if err != nil {
			return nil, err
		}

		temp.Jumlah = count
		temp.Judul = book
		temp.Judul.DetailBuku = detailBuku
		resp = append(resp, temp)
	}

	return resp, nil
}

func (b *bukuController) InputJudulBaru(input models.InputJudulBaruRequest) error {
	var (
		idJudul int
		err     error
	)
	// Input judul baru
	err = b.db.Transaction(func(tx *gorm.DB) error {
		id, txerr := b.bukuDb.InsertJudulBaru(input.DetailJudul, tx)
		idJudul = id
		if txerr != nil {
			return txerr
		}

		if len(input.DetailJudul.IDBuku) > 0 {
			for _, idBuku := range input.DetailJudul.IDBuku {
				txerr = b.bukuDb.InsertIDBuku(idBuku, idJudul, tx)
				if txerr != nil {
					return txerr
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	if len(input.DetailJudul.IDBuku) > 0 {
		return nil
	}

	jumlahBuku := input.Jumlah
	buku := []models.Buku{}

	for i := 0; i < jumlahBuku; i++ {
		temp := models.Buku{}
		temp.Status = util.STATUS_TERSEDIA
		temp.IdJudul = idJudul
		buku = append(buku, temp)
	}

	err = b.db.Transaction(func(tx *gorm.DB) error {
		txerr := b.bukuDb.BatchInsertBukuIndividualBaru(buku, tx)
		if txerr != nil {
			return txerr
		}
		return nil
	})

	if err != nil {
		return err
	}

	bukuElastic, err := mapperJudulToJudulElastic(b, &input, idJudul)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	b.bukuEs.InsertBuku(bukuElastic)

	return nil
}

func (b *bukuController) UpdateJudulBuku(input models.Judul) error {
	var (
		idJudul = input.Id
	)


	oldJudul, err := b.bukuDb.GetJudulByID(idJudul)
	if err != nil {
		return err
	}

	MapAndValidateUpdateJudul(&input, &oldJudul)

	err = b.db.Transaction(func(tx *gorm.DB) error {
		txerr := b.bukuDb.UpdateJudul(input, tx)
		if txerr != nil {
			return txerr
		}

		return nil
	})

	if err != nil {
		return err
	}

	bukuElastic, err := mapperUpdateToElastic(b, &input)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	b.bukuEs.InsertBuku(bukuElastic)

	return nil
}

func (b *bukuController) SyncBukuDBWithES() error {
	var (
		limit       = 100
		offset      = 0
		judulCount  int
		kategoriMap = make(map[int]string)
	)

	query := "SELECT COUNT(*) FROM judul"

	

	err := b.db.Raw(query).Scan(&judulCount).Error
	if err != nil {
		return err
	}

	for offset < judulCount {
		limitOffset := models.FilterBuku{
			Limit:  limit,
			Offset: offset,
		}

		judulTemp, err := b.bukuDb.GetAllJudul(limitOffset)
		if err != nil {
			return err
		}

		for _, judul := range judulTemp {
			var (
				tempEs   models.JudulElastic
				kategori string
				err error
			)

			if val, ok := kategoriMap[judul.IDKategori]; ok {
				kategori = val
			} else {
				kategori, err = b.bukuDb.GetKategoriBuku(judul.IDKategori)

				if err != nil {
					logrus.Errorf("[SyncBukuDBWithES] Error get kategori name : %v", err.Error())
				}
				kategoriMap[judul.IDKategori] = kategori
			}

			totalJudul, err := b.bukuDb.GetBukuCountPerJudul(judul.Id)
			if err != nil {
				logrus.Errorf("[SyncBukuDBWithES] Error get total judul count : %v", err.Error())
			}

			totalTersedia, err := b.bukuDb.GetAvailableJudulByIDJudul(judul.Id)
			if err != nil {
				logrus.Errorf("[SyncBukuDBWithES] Error get total judul available count : %v", err.Error())
			}

			tempEs.Bahasa = judul.Bahasa
			tempEs.Id = judul.Id
			tempEs.Judul = judul.Judul
			tempEs.Kategori = kategori
			tempEs.IDKategori = judul.IDKategori
			tempEs.Penulis = judul.Penulis
			tempEs.Tahun = judul.Tahun
			tempEs.Jenis = judul.Jenis
			tempEs.Penerbit = judul.Penerbit
			tempEs.JumlahTotal = totalJudul
			tempEs.JumlahTersedia = totalTersedia
			tempEs.Tipe = "Buku"

			b.bukuEs.InsertBuku(tempEs)
		}
		time.Sleep(10 * time.Second)
		offset = offset + limit
	}

	return nil
}

func (b *bukuController) SearchBuku(input models.ElasticFilter) ([]models.JudulElastic, int, error) {
	var (
		resp []models.JudulElastic
		err error
		count int
	)
	if input.Jenis == "" {
		respBuku, bukuCount, err := b.bukuEs.SearchBuku(input)
		if err != nil {
			return nil, count, err
		}

		respLaporan, laporanCount, err := b.laporanEs.SearchLaporan(input)
		if err != nil {
			return nil, count, err
		}

		count = bukuCount + laporanCount

		resp = append(resp, respBuku...)
		resp = append(resp, respLaporan...)

	} else if input.Jenis == "jurnal" || input.Jenis == "pkl" || input.Jenis == "skripsi" {
		resp, count, err = b.laporanEs.SearchLaporan(input)
	} else if input.Jenis == "buku" {
		resp, count, err = b.bukuEs.SearchBuku(input)
	}

	if err != nil {
		return nil, count, err
	}
	

	return resp, count, nil
}

func (b *bukuController) GetAllKategori() ([]models.Kategori, error) {
	resp, err := b.bukuDb.GetAllKategori()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *bukuController) GetJudulByBukuID(idBuku int) (models.Judul, error) {
	var (
		resp models.Judul
	)

	idJudul, err := b.bukuDb.GetIDJudulByBukuID(idBuku)
	if err != nil {
		return resp, err
	}

	resp, err = b.bukuDb.GetJudulByID(idJudul)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (b *bukuController) GetCartItemsByNIM(nim string) (models.CartResponse, error) {
	var (
		resp models.CartResponse
	)

	idJudul, err := b.bukuDb.GetCartItemsByNIM(nim)
	if err != nil {
		return resp, err
	}

	for _, id := range idJudul {
		var detailJudul models.Judul
		countJudul, err := b.bukuDb.GetAvailableJudulByIDJudul(id)
		if err != nil {
			return resp, err
		}

		detailJudul, err = b.bukuDb.GetJudulByID(id)
		if err != nil {
			return resp, err
		}

		detailJudul.Count = int64(countJudul)

		if countJudul < 1 {
			detailJudul.IsAvailable = false
		} else {
			detailJudul.IsAvailable = true
		}


		resp.Judul = append(resp.Judul, detailJudul)
	}

	return resp, nil

}

func (b *bukuController) GetJudulByIDElastic(input models.ElasticFilter) (models.JudulElastic, error) {
	resp, err := b.bukuEs.GetJudulByID(input)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (b *bukuController) DeleteBukuOrLaporan(input models.ElasticFilter) (error) {
	var queryJenis string
	var jumlah int
	switch strings.ToLower(input.Jenis) {
	case "buku":
		queryJenis = "buku"
	default:
		queryJenis = "laporan"
	}

	// get count to check if buku still exist in db
	queryCount := fmt.Sprintf("SELECT COUNT(*) from %v WHERE id = %v", queryJenis, input.ID)
	errCount := b.db.Raw(queryCount).Scan(&jumlah).Error
	if errCount != nil {
		return errCount
	}

	err := b.db.Transaction(func(tx *gorm.DB) error {

		if jumlah > 0 {
			if queryJenis == "buku" {
				txerr := b.bukuDb.DeleteBuku(input.ID, tx)
				if txerr != nil {
					return txerr
				}
			} else {
					txerr := b.laporanDB.DeleteLaporan(input.ID, tx)
					if txerr != nil {
						return txerr
					}
			}
		}

		// delete from ES
		txerr := b.bukuEs.DeleteBuku(input)
		if txerr != nil {
			return txerr
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}