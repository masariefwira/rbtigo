package main

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ikalkali/rbti-go/auth"
	cBuku "github.com/ikalkali/rbti-go/controller/buku"
	cLaporan "github.com/ikalkali/rbti-go/controller/laporan"
	cPeminjaman "github.com/ikalkali/rbti-go/controller/peminjaman"
	esBuku "github.com/ikalkali/rbti-go/elastic/buku"
	esLaporan "github.com/ikalkali/rbti-go/elastic/laporan"
	"github.com/ikalkali/rbti-go/entity/buku"
	"github.com/ikalkali/rbti-go/entity/laporan"
	"github.com/ikalkali/rbti-go/entity/mahasiswa"
	"github.com/ikalkali/rbti-go/entity/peminjaman"
	"github.com/ikalkali/rbti-go/handler"
	"github.com/ikalkali/rbti-go/util"
	elastic "github.com/olivere/elastic/v7"
	cron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	err         error
	esClient    *elastic.Client
	cronClient  *cron.Cron
	authService auth.AuthenticationService
)

func init() {
	dsn := util.POSTGRES_DSN
	ctx := context.Background()

	// jakartaTime, _ := time.LoadLocation("Asia/Jakarta")

	log.Info("Init db started...")
	db, err = gorm.Open(postgres.New(
		postgres.Config{
			DSN: dsn,
		},
	), &gorm.Config{})
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("Init DB done")

	// cronClient = cron.New(cron.WithLocation(jakartaTime))
	// log.Info("Init cron done")

	log.Info("Init Elasticsearch started...")
	esClient, err = elastic.NewClient()
	if err != nil {
		panic(err)
	}
	log.Info("Init Elasticsearch done")

	info, code, err := esClient.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
}

func main() {
	r := gin.Default()

	// defer cronClient.Stop()

	peminjamanDb := peminjaman.NewEntity(db)
	bukuDb := buku.NewEntity(db)
	mahasiswaDb := mahasiswa.NewEntity(db)
	laporanDb := laporan.NewEntity(db)
	authService = auth.New(mahasiswaDb, db)

	elasticBuku := esBuku.NewElasticRepo(esClient)
	elasticLaporan := esLaporan.NewElasticRepo(esClient)

	peminjamanController := cPeminjaman.NewController(peminjamanDb, db, bukuDb, mahasiswaDb, cronClient, elasticBuku)
	bukuController := cBuku.NewController(bukuDb, db, elasticBuku, elasticLaporan, laporanDb)
	laporanController := cLaporan.NewController(peminjamanDb, db, bukuDb, mahasiswaDb, elasticLaporan, laporanDb)

	handler := handler.NewHandler(peminjamanController, bukuController, elasticBuku, laporanController, authService)

	// cron
	// cronClient.AddFunc("0 8 * * *", peminjamanController.NotifyPeminjamTelat)
	// cronClient.AddFunc("0 8 * * *", peminjamanController.NotifyPeminjaman)

	// cronClient.Start()
	// log.Print("Cron started...")

	MapURLsAndStartServer(r, handler, authService)
}

func MapURLsAndStartServer(r *gin.Engine, h handler.HandlerInterface, a auth.AuthenticationService) {
	// config := cors.DefaultConfig()
	// config.AllowAllOrigins = true
	// config.AllowMethods = []string{"GET", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"}

	authMiddleware := a.NewMiddleware()

	r.Use(cors.Default())
	r.GET("/ping", h.Ping)
	r.GET("api/buku/sync", h.SyncES)           // Elastic
	r.GET("api/laporan/sync", h.SyncESLaporan) // Elastic
	r.GET("api/kategori", h.GetAllKategori)
	r.GET("api/buku", h.GetJudulByIDBuku)
	r.GET("api/mahasiswa", h.GetMahasiswaByNIM)
	r.GET("api/paper", h.SearchPaper)
	r.GET("api/paper/child", h.GetAllChildFromPaper)

	r.POST("signup", h.Signup)
	r.POST("login", authMiddleware.LoginHandler)
	r.GET("auth/refresh", authMiddleware.RefreshHandler)

	r.POST("api/peminjaman/nim", h.GetDetailPeminjamanByNIM)
	r.POST("api/peminjaman", h.InsertPeminjaman)
	r.POST("api/buku", h.GetAllJudul)
	r.POST("api/buku/search", h.SearchBuku) // Elastic
	r.POST("api/peminjaman/all", h.GetAllPeminjaman)
	r.POST("api/cart/edit", h.EditBukuInCart)
	r.POST("api/cart", h.GetCartItemsByNIM)
	r.POST("api/peminjaman/detail", h.GetDetailPeminjamanByID)
	r.POST("api/buku/judul", h.GetJudulByIDElastic)
	r.POST("api/mahasiswa", h.InputMahasiswaBaru)
	r.POST("api/laporan", h.InsertLaporan)
	r.POST("api/paper", h.InsertPaper)
	r.POST("api/karya_tulis", h.InsertKaryaTulis)
	r.POST("api/paper/all", h.GetAllPaper)
	r.POST("api/karya_tulis/all", h.GetAllArtikel)

	r.PATCH("api/buku", h.InsertJudulBaru)
	r.PATCH("api/peminjaman", h.InsertPengembalian)

	r.PUT("api/buku", h.UpdateJudulBuku)

	r.DELETE("api/buku", h.DeleteBukuOrLaporan)

	r.Run("0.0.0.0:8080")
}
