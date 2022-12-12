package mahasiswa

import (
	"fmt"

	models "github.com/ikalkali/rbti-go/entity/models"
	"gorm.io/gorm"
)

type MahasiswaEntityInterface interface {
	GetAllMahasiswa() []models.Mahasiswa
	GetMahasiswaByNIM(nim string) (models.Mahasiswa, error)
	InsertMahasiswaBaru(input models.Mahasiswa) error
	GetPasswordByEmail(email string) (string, error)
	GetDetailMahasiswaByEmail(email string) (models.Mahasiswa, error)
	InputMahasiswaBaruSignup(input models.Mahasiswa, tx *gorm.DB) error
	UpdatePasswordMahasiswa(input models.Mahasiswa, tx *gorm.DB) error
	CheckMahasiswaExists(nim string) (bool, error)
	CheckMahasiswaExistsEligible(nim string) (bool, error)
}

type mahasiswa struct {
	db *gorm.DB
}

func NewEntity(db *gorm.DB) *mahasiswa {
	return &mahasiswa{db}
}

func (m *mahasiswa) GetAllMahasiswa() []models.Mahasiswa {
	var mhs []models.Mahasiswa
	err := m.db.Find(&mhs).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	return mhs
}

func (m *mahasiswa) GetMahasiswaByNIM(nim string) (models.Mahasiswa, error) {
	var mhs models.Mahasiswa
	err := m.db.Where("nim = ?", nim).First(&mhs).Error
	if err != nil {
		return mhs, err
	}

	return mhs, nil
}

func (m *mahasiswa) InsertMahasiswaBaru(input models.Mahasiswa) error {
	err := m.db.Exec(InsertMahasiswaBaru, input.Nim, input.Nama, input.Nomor_telp, input.Email, input.Angkatan).Error
	if err != nil {
		return err
	}

	return nil
}

func (m *mahasiswa) GetPasswordByEmail(email string) (string, error) {
	var password string
	err := m.db.Raw(GetPasswordByEmail, email).Scan(&password).Error
	if err != nil {
		return password, err
	}

	return password, nil
}

func (m *mahasiswa) GetDetailMahasiswaByEmail(email string) (models.Mahasiswa, error) {
	var response models.Mahasiswa
	err := m.db.Raw(GetDetailMahasiswaByEmail, email).Scan(&response).Error
	if err != nil {
		return response, err
	}

	return response, nil
}

func (m *mahasiswa) InputMahasiswaBaruSignup(input models.Mahasiswa, tx *gorm.DB) error {
	err := tx.Exec(
		InsertMahasiswaBaruSignup,
		input.Nim,
		input.Nama,
		input.Nomor_telp,
		input.Email,
		input.Angkatan,
		input.Password,
		input.Role,
	).Error

	if err != nil {
		return err
	}

	return nil
}

func (m *mahasiswa) UpdatePasswordMahasiswa(input models.Mahasiswa, tx *gorm.DB) error {
	err := tx.Exec(
		UpdatePasswordMahasiswa,
		input.Email,
		input.Password,
		input.Nim,
	).Error

	if err != nil {
		return err
	}

	return nil
}

func (m *mahasiswa) CheckMahasiswaExists(nim string) (bool, error) {
	var count int
	err := m.db.Raw(CheckMahasiswaExists, nim).Scan(&count).Error
	if err != nil {
		return false, nil
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (m *mahasiswa) CheckMahasiswaExistsEligible(nim string) (bool, error) {
	var count int
	err := m.db.Raw(CheckMahasiswaExistEligible, nim).Scan(&count).Error
	if err != nil {
		return false, nil
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
