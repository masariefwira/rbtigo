package peminjaman

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	esBuku "github.com/ikalkali/rbti-go/elastic/buku"
	mailer "github.com/ikalkali/rbti-go/email"
	"github.com/ikalkali/rbti-go/entity/buku"
	"github.com/ikalkali/rbti-go/entity/mahasiswa"
	models "github.com/ikalkali/rbti-go/entity/models"
	"github.com/ikalkali/rbti-go/entity/peminjaman"
	"github.com/ikalkali/rbti-go/util"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PeminjamanControllerInterface interface {
	GetDetailPeminjamanByNIM(nim string) (models.DetailPeminjamanByNIM, error)
	InsertPeminjaman(request models.InputPeminjamanRequest) (int, error)
	InsertPengembalian(request []models.InputPengembalianRequest) error
	GetMahasiswaByNim(nim string) (models.Mahasiswa, error)
	GetAllPeminjaman(filter models.PeminjamanFilter) ([]models.PeminjamanResponse, error)
	EditBukuInCart(input models.CartRequest) (error)
	GetDetailPeminjamanByID(idPeminjaman int) (models.Peminjaman, error)
	InputMahasiswaBaru(input models.Mahasiswa) (error)
}

type peminjamanController struct {
	peminjamanEntity peminjaman.PeminjamanEntityInterface
	bukuEntity       buku.BukuEntityInterface
	mahasiswaEntity  mahasiswa.MahasiswaEntityInterface
	db               *gorm.DB
	cron             *cron.Cron
	esBuku               esBuku.BukuElasticInterface 
}

func NewController(
	entity peminjaman.PeminjamanEntityInterface,
	db *gorm.DB,
	bukuEntity buku.BukuEntityInterface,
	mahasiswaEntity mahasiswa.MahasiswaEntityInterface,
	cron *cron.Cron,
	esBuku esBuku.BukuElasticInterface ,
) *peminjamanController {
	return &peminjamanController{
		peminjamanEntity: entity,
		db:               db,
		bukuEntity:       bukuEntity,
		mahasiswaEntity:  mahasiswaEntity,
		cron:             cron,
		esBuku: esBuku,
	}
}

func (p *peminjamanController) GetDetailPeminjamanByNIM(nim string) (models.DetailPeminjamanByNIM, error) {
	var (
		totalDenda int64
	)

	// Get ID Peminjaman
	detailPeminjaman, _ := p.peminjamanEntity.GetIdPeminjamanByNIM(nim)

	// Populate detail peminjaman
	for idx, pinjam := range detailPeminjaman {
		peminjaman := &detailPeminjaman[idx]
		idPeminjaman := pinjam.IDPeminjaman

		detailBukuRequest := models.DetailBukuPeminjamanFilter{
			IDPeminjaman: idPeminjaman,
			Ketersediaan: util.MASIH_DIPINJAM,
		}

		peminjaman.BukuDipinjam, _ = p.peminjamanEntity.GetJudulDetailBukuDipinjamByIdPeminjaman(detailBukuRequest)
		peminjaman.Denda = int64(util.HitungDenda(peminjaman.TanggalPeminjaman))
		totalDenda = totalDenda + peminjaman.Denda
	}

	resp := models.DetailPeminjamanByNIM{
		Peminjaman: detailPeminjaman,
		TotalDenda: totalDenda,
	}

	return resp, nil
}

func (p *peminjamanController) InsertPeminjaman(request models.InputPeminjamanRequest) (int, error) {
	var (
		status int
		idPeminjaman int
	)

	if request.Source == "app" {
		status = util.STATUS_DIPINJAM_BELUM_DIAMBIL
	} else {
		status = util.STATUS_DIPINJAM
	}

	if request.ID != 0 {
		err := p.db.Transaction(func(tx *gorm.DB) error {
			txerr := p.peminjamanEntity.UpdatePeminjamanKeAktif(request.ID, tx)
			if txerr != nil {
				return txerr
			}

			return nil
		})

		if err != nil {
			return request.ID, err
		}

		return request.ID, nil
	}

	dataPeminjaman := models.PeminjamanDB{
		TanggalPeminjaman:   time.Now(),
		Status:             	status,
		NIM:                 request.NIM,
		TenggatPengembalian: time.Now().Add(util.DURASI_PEMINJAMAN),
	}


	if len(request.IDBuku) == 0 {
		for _, id := range request.IDJudul {
			requestPointer := &request
			idBuku, err := p.bukuEntity.GetAvailableIDBukuByIDJudul(id)
			if err != nil {
				return idPeminjaman, err
			}

			if idBuku == 0 {
				return idPeminjaman, errors.New("[InsertPeminjaman] no books available for the given ID judul")
			}
	
			requestPointer.IDBuku = append(requestPointer.IDBuku, idBuku)
		}
	}

	err := p.db.Transaction(func(tx *gorm.DB) error {
		// Insert detail peminjaman to peminjaman table
		txerr := p.peminjamanEntity.InsertPeminjaman(&dataPeminjaman, tx)
		if txerr != nil {
			return txerr
		}

		idPeminjaman = dataPeminjaman.ID
		dataPeminjamanMap := []models.PeminjamanMap{}
		for _, idBuku := range request.IDBuku {
			peminjamanMap := models.PeminjamanMap{}
			peminjamanMap.IDBuku = idBuku
			peminjamanMap.IDPeminjaman = idPeminjaman
			dataPeminjamanMap = append(dataPeminjamanMap, peminjamanMap)

			// Update status buku dipinjam
			txerr = p.bukuEntity.UpdateBukuDipinjam(idBuku, util.STATUS_DIPINJAM, tx)
			if txerr != nil {
				return txerr
			}

			idJudul, errGet := p.bukuEntity.GetIDJudulByBukuID(idBuku)
			txerr = errGet
			if txerr != nil {
				return txerr
			}

			txerr = p.esBuku.ChangeBukuStock(idJudul, "decrement")
			if txerr != nil {
				return txerr
			}
		}

		// Create new peminjaman_buku map
		txerr = p.peminjamanEntity.BatchInsertPeminjamanMap(dataPeminjamanMap, tx)
		if txerr != nil {
			return txerr
		}

		return nil
	})

	if err != nil {
		return idPeminjaman, err
	}

	return idPeminjaman, err
}

func (p *peminjamanController) InsertPengembalian(request []models.InputPengembalianRequest) error {
	// Init transaction
	err := p.db.Transaction(func(tx *gorm.DB) error {
		for _, peminjaman := range request {
			idPeminjaman := peminjaman.IDPeminjaman
			for _, buku := range peminjaman.IDBuku {

				// Update map
				txerr := p.peminjamanEntity.UpdateBukuKembali(buku, idPeminjaman, p.db)
				if txerr != nil {
					return txerr
				}

				// Update buku
				txerr = p.bukuEntity.UpdateBukuDipinjam(buku, util.STATUS_TERSEDIA, tx)
				if txerr != nil {
					return txerr
				}

				idJudul, errGet := p.bukuEntity.GetIDJudulByBukuID(buku)
				txerr = errGet
				if txerr != nil {
					return txerr
				}

				txerr = p.esBuku.ChangeBukuStock(idJudul, "increment")
				if txerr != nil {
					return txerr
				}

			}
			// Check if all buku on peminjaman is returned
			allKembali, err := p.peminjamanEntity.CheckBukuAllKembali(idPeminjaman)
			if err != nil {
				return err
			}

			// Updates to kembali partial if not all kembali
			if !allKembali {
				txerr := p.peminjamanEntity.UpdatePeminjamanKembali(idPeminjaman, util.STATUS_KEMBALI_PARTIAL, tx)
				if txerr != nil {
					return txerr
				}
				// Updates to kembali all if all returned
			} else {
				txerr := p.peminjamanEntity.UpdatePeminjamanKembali(idPeminjaman, util.STATUS_KEMBALI, tx)
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

	return nil
}

func (p *peminjamanController) GetMahasiswaByNim(nim string) (models.Mahasiswa, error) {
	var (
		resp models.Mahasiswa
	)

	resp, err := p.mahasiswaEntity.GetMahasiswaByNIM(nim)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (p *peminjamanController) GetAllPeminjaman(filter models.PeminjamanFilter) ([]models.PeminjamanResponse, error) {
	var (
		resp         []models.PeminjamanResponse
		mapMahasiswa = make(map[string]models.Mahasiswa)
	)

	resp, err := p.peminjamanEntity.GetAllPeminjaman(filter)
	if err != nil {
		return resp, err
	}

	for idx, _ := range resp {
		temp := &resp[idx]
		nim := temp.NIM

		if val, ok := mapMahasiswa[nim]; ok {
			temp.Nama = val.Nama
			temp.Email = val.Email
		} else {
			mhs, err := p.mahasiswaEntity.GetMahasiswaByNIM(nim)
			if err != nil {
				return resp, err
			}

			nama := mhs.Nama
			mapMahasiswa[nim] = mhs

			temp.Nama = nama
			temp.Email = mhs.Email
		}
		days := math.Round(time.Since(temp.TenggatPengembalian).Hours() / 24)
		if days < 0 {
			temp.Keterlambatan = ""
		} else {
			if temp.Status != util.STATUS_KEMBALI {
				temp.Keterlambatan = fmt.Sprintf("%v hari", days)
			} else {
				temp.Keterlambatan = ""
			}
		}

	}

	return resp, nil
}

func (p *peminjamanController) NotifyPeminjamTelat() {
	filter := models.PeminjamanFilter{
		Status: []int{2, 5},
		Telat:  "true",
		Limit:  20,
	}

	resp, err := p.GetAllPeminjaman(filter)
	if err != nil {
		logrus.Errorf("[NotifyPeminjamTelat] fail to GetAllPeminjaman : %v", err)
		return
	}

	for _, data := range resp {
		detailBukuRequest := models.DetailBukuPeminjamanFilter{
			IDPeminjaman: data.ID,
			Ketersediaan: util.MASIH_DIPINJAM,
		}
		data.DetailPeminjaman, _ = p.peminjamanEntity.GetJudulDetailBukuDipinjamByIdPeminjaman(detailBukuRequest)
		data.Denda = int64(util.HitungDenda(data.TanggalPeminjaman))

		msg := mailer.MapEmailKeterlambatan(data)
		err := mailer.SendToMail(data.Email, "[Reminder] Tenggat Pengembalian Buku", msg)
		if err != nil {
			log.Printf("[NotifyPeminjamTelat] SendToMail error : %v", err)
		}
		log.Print("[NotifyPeminjamTelat] Mail sent")
	}

	log.Print("All email sent")

}

func (p *peminjamanController) NotifyPeminjaman() {
	filter := models.PeminjamanFilter{
		Status: []int{2, 5, 4},
		Limit:  100,
		Waktu: models.PeminjamanFilterWaktu{
			Waktu:    "3",
			Operator: "=",
		},
	}

	resp, err := p.GetAllPeminjaman(filter)
	if err != nil {
		logrus.Errorf("[NotifyPeminjaman] fail to GetAllPeminjaman : %v", err)
		return
	}

	fmt.Printf("\n%+v\n", resp)

	for _, data := range resp {
		detailBukuRequest := models.DetailBukuPeminjamanFilter{
			IDPeminjaman: data.ID,
			Ketersediaan: util.MASIH_DIPINJAM,
		}
		data.DetailPeminjaman, _ = p.peminjamanEntity.GetJudulDetailBukuDipinjamByIdPeminjaman(detailBukuRequest)

		msg := mailer.MapEmailReminder(data)
		err := mailer.SendToMail(data.Email, "[Reminder] Pengembalian Buku", msg)
		if err != nil {
			log.Printf("[NotifyPeminjaman] SendToMail error : %v", err)
		}
		log.Print("[NotifyPeminjaman] Mail sent")
	}
}

func (p *peminjamanController) EditBukuInCart(input models.CartRequest) (error) {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		for _, idJudul := range input.IdJudul {

			var txerr error

			switch input.Action {
			case "add":
				txerr = p.bukuEntity.AddBukuToCart(idJudul, input.NIM, tx)
			case "remove":
				txerr = p.bukuEntity.DeleteItemFromCart(idJudul, input.NIM, tx)
			default:
				txerr = errors.New("[EditBukuInCart] invalid action given, the available actions are 'add' and 'remove'")
			}

			if txerr != nil {
				return txerr
			}

		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *peminjamanController) GetDetailPeminjamanByID(idPeminjaman int) (models.Peminjaman, error) {
	var (
		resp models.Peminjaman
		err error
	)

	resp, err = p.peminjamanEntity.GetDetailPeminjamanByID(idPeminjaman)
	if err != nil {
		return resp, err
	}

	if resp.Status == 0 {
		return resp, errors.New("id peminjaman tidak ada!")
	}

	inputDetailBuku := models.DetailBukuPeminjamanFilter{
		IDPeminjaman: idPeminjaman,
		Ketersediaan: "semua",
	}

	resp.BukuDipinjam, err = p.peminjamanEntity.GetJudulDetailBukuDipinjamByIdPeminjaman(inputDetailBuku)
	if err != nil {
		return resp, err
	}

	resp.Denda = int64(util.HitungDendaByTenggat(resp.TenggatPengembalian))
	resp.IDPeminjaman = idPeminjaman

	return resp, nil
}

func (p *peminjamanController) InputMahasiswaBaru(input models.Mahasiswa) (error) {
	if input.Nama == "" {
		return errors.New("[InputMahasiswaBaru] nama tidak boleh kosong")
	}

	err := p.mahasiswaEntity.InsertMahasiswaBaru(input)
	if err != nil {
		return err
	}

	return nil
}