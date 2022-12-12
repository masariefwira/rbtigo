package laporan

import (
	"fmt"
	"strings"

	esLaporan "github.com/ikalkali/rbti-go/elastic/laporan"
	"github.com/ikalkali/rbti-go/entity/buku"
	"github.com/ikalkali/rbti-go/entity/laporan"
	"github.com/ikalkali/rbti-go/entity/mahasiswa"
	"github.com/ikalkali/rbti-go/entity/models"
	"github.com/ikalkali/rbti-go/entity/peminjaman"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LaporanControllerInterface interface {
	SyncESLaporan() error
	InsertLaporan(input models.LaporanDB) error
	InsertPaper(input models.Paper) error
	SearchPaper(query string, jenis int) ([]models.Paper, int, error)
	InsertKaryaTulis(input models.KaryaTulis) error
	GetAllChildFromPaper(idPaper int) ([]models.JudulElastic, error)
	GetAllPaper(filter models.PaperFilter) ([]models.Paper, int, error)
	GetAllArtikel(filter models.PaperFilter) ([]models.Artikel, int, error)
}

type laporanController struct {
	peminjamanEntity peminjaman.PeminjamanEntityInterface
	bukuEntity       buku.BukuEntityInterface
	mahasiswaEntity  mahasiswa.MahasiswaEntityInterface
	laporanEntity    laporan.LaporanEntityInterface
	db               *gorm.DB
	esLaporan        esLaporan.LaporanElasticInterface
}

func NewController(
	entity peminjaman.PeminjamanEntityInterface,
	db *gorm.DB,
	bukuEntity buku.BukuEntityInterface,
	mahasiswaEntity mahasiswa.MahasiswaEntityInterface,
	esLaporan esLaporan.LaporanElasticInterface,
	laporanEntity laporan.LaporanEntityInterface,
) *laporanController {
	return &laporanController{
		peminjamanEntity: entity,
		db:               db,
		bukuEntity:       bukuEntity,
		mahasiswaEntity:  mahasiswaEntity,
		esLaporan:        esLaporan,
		laporanEntity:    laporanEntity,
	}
}

func (l *laporanController) SyncESLaporan() error {
	var (
		limit       = 100
		offset      = 0
		mapKategori = make(map[int]string)
	)

	judulCount, err := l.laporanEntity.GetLaporanCount()
	if err != nil {
		logrus.Error(err)
		return err
	}

	for offset < judulCount {
		filterLaporan := models.LaporanFilter{
			Offset: offset,
			Limit:  limit,
		}

		judulTemp, err := l.laporanEntity.GetAllLaporan(filterLaporan)
		if err != nil {
			return err
		}

		for _, judul := range judulTemp {
			var (
				tempEs   models.JudulElastic
				kategori string
			)

			if val, ok := mapKategori[judul.IdKategori]; ok {
				kategori = val
			} else {
				kategori, err = l.bukuEntity.GetKategoriBuku(judul.IdKategori)
				if err != nil {
					return err
				}

				mapKategori[judul.IdKategori] = kategori
			}

			penulis, err := l.mahasiswaEntity.GetMahasiswaByNIM(judul.NIM)
			if err != nil {
				return err
			}

			tempEs.Id = judul.Id
			tempEs.Tipe = judul.Jenis
			tempEs.Tahun = judul.Tahun
			tempEs.NIM = judul.NIM
			tempEs.Penulis = strings.Title(strings.ToLower(penulis.Nama))
			tempEs.Kategori = kategori
			tempEs.IDKategori = judul.IdKategori
			tempEs.Judul = strings.Title(strings.ToLower(judul.Judul))

			fmt.Printf("%+v\n\n", tempEs)

			err = l.esLaporan.InsertLaporan(tempEs)
			if err != nil {
				logrus.Error(err)
			}
		}

		offset = limit + offset
	}
	logrus.Info("[SyncLaporanES] sync laporan to elasticsearch done")
	return nil
}

func (l *laporanController) InsertLaporan(input models.LaporanDB) error {
	var idLaporan int
	err := l.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if input.IdPaper != 0 {
			idLaporan, err = l.laporanEntity.InsertArtikelOrMakalah(input, tx)
		} else {
			idLaporan, err = l.laporanEntity.InsertLaporan(input, tx)
		}

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	penulis, err := l.mahasiswaEntity.GetMahasiswaByNIM(input.NIM)
	if err != nil {
		logrus.Error(err)
	}

	kategori, err := l.bukuEntity.GetKategoriBuku(input.IdKategori)
	if err != nil {
		logrus.Error(err)
	}

	elasticLaporan := models.JudulElastic{
		Id:         idLaporan,
		Judul:      input.Judul,
		Tahun:      input.Tahun,
		NIM:        input.NIM,
		IDKategori: input.IdKategori,
		Penulis:    strings.Title(strings.ToLower(penulis.Nama)),
		Kategori:   kategori,
		Tipe:       input.Jenis,
	}

	err = l.esLaporan.InsertLaporan(elasticLaporan)
	if err != nil {
		logrus.Error(err)
	}

	return nil
}

func (l *laporanController) InsertPaper(input models.Paper) error {
	err := l.db.Transaction(func(tx *gorm.DB) error {
		_, txerr := l.laporanEntity.InsertPaper(input, tx)
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

func (l *laporanController) SearchPaper(query string, jenis int) ([]models.Paper, int, error) {
	paper, count, err := l.laporanEntity.SearchPaper(query, jenis)
	if err != nil {
		return paper, count, err
	}

	return paper, count, nil
}

func (l *laporanController) GetAllPaper(filter models.PaperFilter) ([]models.Paper, int, error) {
	paper, count, err := l.laporanEntity.GetAllPaper(filter)
	if err != nil {
		return paper, count, err
	}

	return paper, count, nil
}

func (l *laporanController) GetAllArtikel(filter models.PaperFilter) ([]models.Artikel, int, error) {
	// get all artikel
	artikel, count, err := l.laporanEntity.GetAllArtikel(filter)
	if err != nil {
		return artikel, count, err
	}

	fmt.Printf("ARTIKEL %+v", artikel)

	// populate every artikel with judul induk
	for idx, artikelItem := range artikel {
		artikelSelected := &artikel[idx]
		fmt.Printf("ARTIKEL SELECTED %+v", artikelSelected)

		idArtikel := artikelItem.IdPaper
		paper, err := l.laporanEntity.GetDetailPaperByID(idArtikel)
		if err != nil {
			return artikel, count, err
		}

		judulGabung := fmt.Sprintf("%v, %v", paper.Judul, paper.Volume)
		artikelSelected.JudulInduk = judulGabung
	}

	fmt.Printf("ARTIKEL AFTER %+v", artikel)

	return artikel, count, nil
}

func (l *laporanController) InsertKaryaTulis(input models.KaryaTulis) error {
	l.db.Transaction(func(tx *gorm.DB) error {
		id, txerr := l.laporanEntity.InsertKaryaTulis(input, tx)
		if txerr != nil {
			return txerr
		}

		input.Id = id
		judulElastic := mapKaryaTulisToElastic(input)

		kategori, txerr := l.bukuEntity.GetKategoriBuku(input.IDKategori)
		if txerr != nil {
			return txerr
		}

		judulElastic.Kategori = kategori

		fmt.Printf("JUDUL ELASTIC INSERT %v", judulElastic)

		txerr = l.esLaporan.InsertKaryaTulis(judulElastic)
		if txerr != nil {
			return txerr
		}

		return nil
	})

	return nil
}

func mapKaryaTulisToElastic(input models.KaryaTulis) models.JudulElastic {
	judulElastic := models.JudulElastic{}

	judulElastic.Id = input.Id
	judulElastic.Judul = input.Judul
	judulElastic.Penulis = input.Penulis
	judulElastic.IDKategori = input.IDKategori
	judulElastic.Tahun = input.Tahun
	judulElastic.IdPaper = input.IDPaper

	return judulElastic
}

func (l *laporanController) GetAllChildFromPaper(idPaper int) ([]models.JudulElastic, error) {
	resp := []models.JudulElastic{}

	// get all ID Artikel from ID Paper
	idKaryaTulis, err := l.laporanEntity.GetIDKaryaTulisByIDPaper(idPaper)
	if err != nil {
		return resp, err
	}

	fmt.Printf("ID Karya Tulis : %v", idKaryaTulis)

	for _, id := range idKaryaTulis {
		temp, err := l.esLaporan.GetKaryaTulisByID(id)
		if err != nil {
			return resp, err
		}

		fmt.Printf("Single Artikel Elastic : %v", temp)

		resp = append(resp, temp)
	}

	return resp, nil
}
