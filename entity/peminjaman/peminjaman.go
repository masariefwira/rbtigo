package peminjaman

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	models "github.com/ikalkali/rbti-go/entity/models"
	"github.com/ikalkali/rbti-go/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PeminjamanEntityInterface interface {
	GetJudulDetailBukuDipinjamByIdPeminjaman(input models.DetailBukuPeminjamanFilter) ([]models.DetailPeminjaman, error)
	GetIdPeminjamanByNIM(nim string) ([]models.Peminjaman, error)
	InsertPeminjaman(input *models.PeminjamanDB, tx *gorm.DB) error
	BatchInsertPeminjamanMap(input []models.PeminjamanMap, tx *gorm.DB) error
	UpdateBukuKembali(idBuku int, idPeminjaman int, tx *gorm.DB) error
	CheckBukuAllKembali(idPeminjaman int) (bool, error)
	UpdatePeminjamanKembali(idPeminjaman int, status int, tx *gorm.DB) error
	GetAllPeminjaman(filter models.PeminjamanFilter) ([]models.PeminjamanResponse, error)
	GetDetailPeminjamanByID(id int) (models.Peminjaman, error)
	UpdatePeminjamanKeAktif(idPeminjaman int, tx *gorm.DB) (error)
}

type peminjaman struct {
	db *gorm.DB
}

func NewEntity(db *gorm.DB) *peminjaman {
	return &peminjaman{db}
}

func (p *peminjaman) GetJudulDetailBukuDipinjamByIdPeminjaman(input models.DetailBukuPeminjamanFilter) ([]models.DetailPeminjaman, error) {
	var (
		detailBuku []models.DetailPeminjaman
		idBuku     int
		judul      string
		tahun      int
		penerbit   string
		query string
	)

	if input.Ketersediaan == util.MASIH_DIPINJAM {
		query = GetDetailBukuDipinjamByIdPeminjaman
	} else if input.Ketersediaan == util.SEMUA_BUKU {
		query = GetDetailBukuDipinjamByIdPeminjamanAll
	}

	rows, err := p.db.Raw(query, input.IDPeminjaman).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		detail := models.DetailPeminjaman{}
		rows.Scan(&idBuku, &judul, &tahun, &penerbit)
		detail.IDBuku = idBuku
		detail.Judul = judul
		detail.Penerbit = penerbit
		detail.Tahun = tahun
		detailBuku = append(detailBuku, detail)
	}

	return detailBuku, nil

}

func (p *peminjaman) GetIdPeminjamanByNIM(nim string) ([]models.Peminjaman, error) {
	var (
		idPeminjaman      int
		tanggalPeminjaman time.Time
		peminjaman        []models.Peminjaman
	)

	rows, err := p.db.Raw(GetIdPeminjamanByNIM, nim).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		pinjam := models.Peminjaman{}
		rows.Scan(&idPeminjaman, &tanggalPeminjaman, &pinjam.Status, &pinjam.TenggatPengembalian)
		pinjam.IDPeminjaman = idPeminjaman
		pinjam.TanggalPeminjaman = tanggalPeminjaman
		peminjaman = append(peminjaman, pinjam)
	}

	return peminjaman, nil
}

func (p *peminjaman) InsertPeminjaman(input *models.PeminjamanDB, tx *gorm.DB) error {
	if err := tx.Create(&input).Error; err != nil {
		return err
	}

	return nil
}

func (p *peminjaman) BatchInsertPeminjamanMap(input []models.PeminjamanMap, tx *gorm.DB) error {
	if err := tx.Create(&input).Error; err != nil {
		return err
	}

	return nil
}

func (p *peminjaman) UpdateBukuKembali(idBuku int, idPeminjaman int, tx *gorm.DB) error {
	now := time.Now()
	if err := tx.Exec(UpdateMapKembali, now, idBuku, idPeminjaman).Error; err != nil {
		return err
	}

	return nil
}

func (p *peminjaman) CheckBukuAllKembali(idPeminjaman int) (bool, error) {
	var (
		jumlah int
	)
	err := p.db.Raw(CheckBukuAllKembali, idPeminjaman).Scan(&jumlah).Error
	if err != nil {
		return false, err
	}

	if jumlah > 0 {
		return false, nil
	}

	return true, nil
}

func (p *peminjaman) UpdatePeminjamanKembali(idPeminjaman int, status int, tx *gorm.DB) error {
	if err := tx.Exec(UpdatePeminjamanKembali, status, idPeminjaman).Error; err != nil {
		return err
	}

	return nil
}

func (p *peminjaman) GetAllPeminjaman(filter models.PeminjamanFilter) ([]models.PeminjamanResponse, error) {
	var (
		limit     = filter.Limit
		offset    = filter.Offset
		waktu     = filter.Waktu
		tempQuery string
		resp      []models.PeminjamanResponse
		count     int64
	)

	isTelat, err := strconv.ParseBool(filter.Telat)
	if err != nil {
		logrus.Errorf("[GetAllPeminjaman] error parsing bool %v", err.Error())
	}

	if len(filter.Status) > 0 {
		trimString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(filter.Status)), ","), "[]")
		tempQuery = tempQuery + fmt.Sprintf(" where status in (%v)", trimString)
	}

	if isTelat {
		if len(filter.Status) > 0 {
			tempQuery = tempQuery + " and tenggat_pengembalian < now()"
		} else {
			tempQuery = tempQuery + " where tenggat_pengembalian < now()"
		}
	} else {
		if waktu.Waktu != "" {
			if waktu.Operator == "=" {
				if len(filter.Status) > 0 {
					tempQuery = tempQuery + fmt.Sprintf(" and current_date + interval '%v days' = tenggat_pengembalian::date", waktu.Waktu)
				} else {
					tempQuery = tempQuery + fmt.Sprintf(" where current_date + interval '%v days' = tenggat_pengembalian::date", waktu.Waktu)
				}
			} else {
				if len(filter.Status) > 0 {
					tempQuery = tempQuery + " and tanggal_peminjaman > now() - interval " + fmt.Sprintf("'%v days'", waktu.Waktu)
				} else {
					tempQuery = tempQuery + " where tanggal_peminjaman > now() - interval " + fmt.Sprintf("'%v days'", waktu.Waktu)
				}
			}
		}
	}

	limitQuery := fmt.Sprintf(" ORDER BY tanggal_peminjaman asc LIMIT %v OFFSET %v", limit, offset)
	finalQuery := GetAllPeminjaman + tempQuery + limitQuery
	fmt.Println(finalQuery)

	rows, err := p.db.Raw(finalQuery).Rows()
	if err != nil {
		return resp, err
	}

	countQ := "SELECT COUNT(*) from peminjaman"
	countQuery := countQ + tempQuery

	fmt.Println(countQuery)

	err = p.db.Raw(countQ).Scan(&count).Error
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		temp := models.PeminjamanResponse{}
		err = rows.Scan(
			&temp.ID,
			&temp.TanggalPeminjaman,
			&temp.Status,
			&temp.NIM,
			&temp.TenggatPengembalian,
		)
		if err != nil {
			return resp, err
		}

		temp.Count = count

		resp = append(resp, temp)
	}

	return resp, nil
}

func (p *peminjaman) GetDetailPeminjamanByID(id int) (models.Peminjaman, error) {
	var (
		resp models.Peminjaman
	)


	rows, err := p.db.Raw(GetDetailPeminjamanByID, id).Rows()
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		rows.Scan(&resp.TanggalPeminjaman, &resp.Status, &resp.NIM, &resp.TenggatPengembalian)
	}

	fmt.Printf("RESP DB %+v\n", resp)

	return resp, nil
}

func (p *peminjaman) UpdatePeminjamanKeAktif(idPeminjaman int, tx *gorm.DB) (error) {
	err := tx.Exec(UpdatePeminjamanKembali, util.STATUS_DIPINJAM, idPeminjaman).Error
	if err != nil {
		return err
	}

	return nil
}