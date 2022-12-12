package auth

import (
	"errors"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/ikalkali/rbti-go/entity/mahasiswa"
	"github.com/ikalkali/rbti-go/entity/models"
	bcrypt "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


type AuthenticationService interface {
	NewMiddleware() *jwt.GinJWTMiddleware
	Authenticate(c *gin.Context) (interface{}, error)
	Signup(input models.Mahasiswa) (error)
	SetPayloadData(data interface{}) jwt.MapClaims
}

type authenticationService struct {
	mahasiswaDb mahasiswa.MahasiswaEntityInterface
	db     *gorm.DB
}

func New(
	mahasiswaDb mahasiswa.MahasiswaEntityInterface,
	db *gorm.DB,
) (*authenticationService) {
	return &authenticationService{mahasiswaDb: mahasiswaDb, db: db}
}

func (a *authenticationService) NewMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm : "staging",
		Key: []byte("SECRET KEY"),
		Timeout: 3600 * time.Hour,
		MaxRefresh: 3600 * time.Hour,
		IdentityKey: "nim",
		PayloadFunc: a.SetPayloadData,
		IdentityHandler: a.IdentityHandler,
		Authenticator: a.Authenticate,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}

func (a *authenticationService) Authenticate(c *gin.Context) (interface{}, error) {
	var loginVals User
	if err := c.BindJSON(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	email := loginVals.Email
	password := loginVals.Password

	storedPassword, err := a.mahasiswaDb.GetPasswordByEmail(email)
	if err != nil {
		return "", err
	}

	valError := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if valError != nil {
		return "", jwt.ErrFailedAuthentication
	}

	detailMahasiswa, err := a.mahasiswaDb.GetDetailMahasiswaByEmail(email)
	if err != nil {
		return "", err
	}

	return &detailMahasiswa, nil 
}	

func (a *authenticationService) IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
			return &models.Mahasiswa{
				Nim: claims["nim"].(string),
			}
}

func (a *authenticationService) SetPayloadData(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.Mahasiswa); ok {
		return jwt.MapClaims{
			"nim": v.Nim,
			"role" : v.Role,
		}
	}
	return jwt.MapClaims{}
}

func (a *authenticationService) Signup(input models.Mahasiswa) (error) {

	// check if nim already exist but password is null
	exist, err := a.mahasiswaDb.CheckMahasiswaExistsEligible(input.Nim)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("user already exists!")
	}

	// if not exist check if NIM is already in
	exist, err = a.mahasiswaDb.CheckMahasiswaExists(input.Nim)
	if err != nil {
		return err
	}

	// hash cleantext password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return err
	}

	input.Password = string(hashedPassword)

	if exist {
		err = a.db.Transaction(func(tx *gorm.DB) error {
			txerr := a.mahasiswaDb.UpdatePasswordMahasiswa(input, tx)
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

	err = a.db.Transaction(func(tx *gorm.DB) error {
		txerr := a.mahasiswaDb.InputMahasiswaBaruSignup(input, tx)
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
	
