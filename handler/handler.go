package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ikalkali/rbti-go/auth"
	cBuku "github.com/ikalkali/rbti-go/controller/buku"
	cLaporan "github.com/ikalkali/rbti-go/controller/laporan"
	cPeminjaman "github.com/ikalkali/rbti-go/controller/peminjaman"
	esBuku "github.com/ikalkali/rbti-go/elastic/buku"
	models "github.com/ikalkali/rbti-go/entity/models"
	"github.com/ikalkali/rbti-go/util"
	"github.com/sirupsen/logrus"
)

type HandlerInterface interface {
	Ping(c *gin.Context)

	// DB
	GetDetailPeminjamanByNIM(c *gin.Context)
	InsertPeminjaman(c *gin.Context)
	InsertPengembalian(c *gin.Context)
	GetAllJudul(c *gin.Context)
	InsertJudulBaru(c *gin.Context)
	UpdateJudulBuku(c *gin.Context)
	GetAllKategori(c *gin.Context)
	GetJudulByIDBuku(c *gin.Context)
	GetMahasiswaByNIM(c *gin.Context)
	GetAllPeminjaman(c *gin.Context)
	GetCartItemsByNIM(c *gin.Context)
	EditBukuInCart(c *gin.Context)
	GetDetailPeminjamanByID(c *gin.Context)
	InputMahasiswaBaru(c *gin.Context)
	InsertPaper(c *gin.Context)
	SearchPaper(c *gin.Context)
	InsertKaryaTulis(c *gin.Context)
	GetAllPaper(c *gin.Context)
	GetAllArtikel(c *gin.Context)

	// auth
	Signup(c *gin.Context)

	// Elastic
	SyncES(c *gin.Context)
	SearchBuku(c *gin.Context)
	SyncESLaporan(c *gin.Context)
	GetJudulByIDElastic(c *gin.Context)
	InsertLaporan(c *gin.Context)
	DeleteBukuOrLaporan(c *gin.Context)
	GetAllChildFromPaper(c *gin.Context)
}

type handler struct {
	peminjamanController cPeminjaman.PeminjamanControllerInterface
	bukuController       cBuku.BukuControllerInterface
	esBuku               esBuku.BukuElasticInterface
	laporanController    cLaporan.LaporanControllerInterface
	auth                 auth.AuthenticationService
}

func NewHandler(
	peminjaman cPeminjaman.PeminjamanControllerInterface,
	buku cBuku.BukuControllerInterface,
	esBuku esBuku.BukuElasticInterface,
	laporan cLaporan.LaporanControllerInterface,
	auth auth.AuthenticationService,
) HandlerInterface {
	return &handler{
		peminjamanController: peminjaman,
		bukuController:       buku,
		esBuku:               esBuku,
		laporanController:    laporan,
		auth:                 auth,
	}
}

func (h *handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *handler) GetDetailPeminjamanByNIM(c *gin.Context) {
	var request models.DetailPeminjamanByNIMRequest
	err := c.BindJSON(&request)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	// get all peminjaman except returned
	resp, err := h.peminjamanController.GetDetailPeminjamanByNIM(request.NIM)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	c.JSON(200, resp)
}

func (h *handler) InsertPeminjaman(c *gin.Context) {
	var request models.InputPeminjamanRequest
	var resp Response
	err := c.BindJSON(&request)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	idPeminjaman, err := h.peminjamanController.InsertPeminjaman(request)
	if err != nil {
		logrus.Error(err.Error())
		resp.Errors = []string{err.Error()}
		c.JSON(500, resp)
		return
	}

	resp.Data = fmt.Sprint(idPeminjaman)

	c.JSON(200, resp)
}

func (h *handler) InsertPengembalian(c *gin.Context) {
	var request models.InputPengembalianData
	err := c.BindJSON(&request)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	err = h.peminjamanController.InsertPengembalian(request.Data)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	c.JSON(200, "success")
}

func (h *handler) GetAllJudul(c *gin.Context) {
	var request models.FilterBuku
	var resp models.GetAllJudulResponse

	err := c.BindJSON(&request)
	if err != nil {
		resp.Errors = []string{err.Error()}
		c.JSON(500, resp)
		return
	}

	resp.Data, err = h.bukuController.GetAllJudul(request)
	if err != nil {
		resp.Errors = []string{err.Error()}
		c.JSON(500, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *handler) InsertJudulBaru(c *gin.Context) {
	var (
		request models.InputJudulBaruRequest
		resp    models.Response
	)

	err := c.BindJSON(&request)
	fmt.Printf("%+v", request)
	if err != nil {
		resp.Errors = []string{err.Error()}
		c.JSON(500, resp)
		return
	}

	// if request.DetailJudul.Filename != nil {
	// 	if *request.DetailJudul.Filename == "" {
	// 		request.DetailJudul.Filename = nil
	// 	}
	// }

	// if request.DetailJudul.Foto != nil {
	// 	if *request.DetailJudul.Foto == "" {
	// 		request.DetailJudul.Foto = nil
	// 	}
	// }

	err = h.bukuController.InputJudulBaru(request)
	if err != nil {
		resp.Errors = []string{err.Error()}
		c.JSON(500, resp)
		return
	}

	resp.Data = "success"

	c.JSON(201, resp)
}

func (h *handler) UpdateJudulBuku(c *gin.Context) {
	var (
		request models.UpdateJudulBukuRequest
		input   models.Judul
		resp    models.Response
	)

	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	input = request.DetailJudul
	input.Id = request.Id

	err = h.bukuController.UpdateJudulBuku(input)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp.Data = "success"
	c.JSON(201, resp)
}

func (h *handler) SyncES(c *gin.Context) {
	go func() {
		err := h.bukuController.SyncBukuDBWithES()
		if err != nil {
			logrus.Errorf("ERROR WHILE SYNC, : %v", err.Error())
		}
	}()

	resp := models.Response{
		Data: "sync called, please monitor in log",
	}

	c.JSON(201, resp)
}

func (h *handler) SearchBuku(c *gin.Context) {
	var query models.ElasticFilter

	c.BindJSON(&query)
	fmt.Printf("QUERY %+v", query)

	search, count, err := h.bukuController.SearchBuku(query)

	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.SearchBukuResponse{
		Data:  search,
		Count: count,
	}

	c.JSON(200, resp)
}

func (h *handler) GetAllKategori(c *gin.Context) {
	resp, err := h.bukuController.GetAllKategori()
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	c.JSON(200, resp)
}

func (h *handler) GetJudulByIDBuku(c *gin.Context) {
	params := c.Request.URL.Query()
	fmt.Printf("%+v", params)

	idBuku, err := strconv.ParseInt(c.Query("idBuku"), 10, 64)
	resp, err := h.bukuController.GetJudulByBukuID(int(idBuku))
	if err != nil {
		util.ErrorWrapper(err, 404, c)
		return
	}

	c.JSON(200, resp)
}

func (h *handler) GetMahasiswaByNIM(c *gin.Context) {
	nim := c.Query("nim")
	resp, err := h.peminjamanController.GetMahasiswaByNim(nim)
	if err != nil {
		util.ErrorWrapper(err, 404, c)
		return
	}

	respProper := models.ResponseTest{
		Data: resp,
	}
	c.JSON(200, respProper)
}

func (h *handler) GetAllPeminjaman(c *gin.Context) {
	var (
		request models.PeminjamanFilter
	)

	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	data, err := h.peminjamanController.GetAllPeminjaman(request)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.ResponseTest{
		Data: data,
	}

	c.JSON(200, resp)
}

func (h *handler) GetCartItemsByNIM(c *gin.Context) {
	var (
		request models.CartRequest
	)

	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	respCart, err := h.bukuController.GetCartItemsByNIM(request.NIM)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.ResponseTest{
		Data: respCart,
	}

	c.JSON(200, resp)
}

func (h *handler) EditBukuInCart(c *gin.Context) {
	var (
		request models.CartRequest
	)

	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	err = h.peminjamanController.EditBukuInCart(request)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.Response{
		Data: "success",
	}

	c.JSON(201, resp)
}

func (h *handler) SyncESLaporan(c *gin.Context) {
	err := h.laporanController.SyncESLaporan()
	if err != nil {
		util.ErrorWrapper(err, 500, c)
	}

	resp := models.ResponseTest{
		Data: "success",
	}

	c.JSON(201, resp)
}

func (h *handler) GetDetailPeminjamanByID(c *gin.Context) {
	var (
		request models.DetailPeminjamanRequest
	)

	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	detailPeminjaman, err := h.peminjamanController.GetDetailPeminjamanByID(request.IDPeminjaman)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.ResponseTest{
		Data: detailPeminjaman,
	}

	c.JSON(200, resp)
}

func (h *handler) GetJudulByIDElastic(c *gin.Context) {
	var (
		request models.ElasticFilter
	)

	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	data, err := h.bukuController.GetJudulByIDElastic(request)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.ResponseTest{
		Data: data,
	}

	c.JSON(200, resp)
}

func (h *handler) InputMahasiswaBaru(c *gin.Context) {
	var (
		request models.Mahasiswa
	)

	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	err = h.peminjamanController.InputMahasiswaBaru(request)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.ResponseTest{
		Data: "success",
	}

	c.JSON(200, resp)
}

func (h *handler) InsertLaporan(c *gin.Context) {
	var (
		request models.LaporanDB
	)

	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	err = h.laporanController.InsertLaporan(request)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.ResponseTest{
		Data: "success",
	}

	c.JSON(200, resp)
}

func (h *handler) DeleteBukuOrLaporan(c *gin.Context) {
	var (
		request models.ElasticFilter
	)

	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	err = h.bukuController.DeleteBukuOrLaporan(request)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.ResponseTest{
		Data: "success",
	}

	c.JSON(200, resp)
}

func (h *handler) Signup(c *gin.Context) {
	var request models.Mahasiswa

	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	err = h.auth.Signup(request)
	if err != nil {
		util.ErrorWrapper(err, 403, c)
		return
	}

	resp := models.ResponseTest{
		Data: "success",
	}

	c.JSON(200, resp)
}

func (h *handler) InsertPaper(c *gin.Context) {
	var request models.Paper

	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	err = h.laporanController.InsertPaper(request)
	if err != nil {
		util.ErrorWrapper(err, 403, c)
		return
	}

	resp := models.ResponseTest{
		Data: "success",
	}

	c.JSON(200, resp)
}

func (h *handler) SearchPaper(c *gin.Context) {
	query := c.Query("query")
	jenis, _ := strconv.Atoi(c.Query("jenis"))
	fmt.Printf("JENIS %v", jenis)
	paper, count, err := h.laporanController.SearchPaper(query, jenis)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.ResponseTest{
		Data:  paper,
		Count: count,
	}

	c.JSON(200, resp)

}

func (h *handler) InsertKaryaTulis(c *gin.Context) {
	var reqeust models.KaryaTulis
	err := c.BindJSON(&reqeust)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	err = h.laporanController.InsertKaryaTulis(reqeust)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.ResponseTest{
		Data: "success",
	}

	c.JSON(200, resp)
}

func (h *handler) GetAllChildFromPaper(c *gin.Context) {
	idPaper := c.Query("idPaper")
	idPaperConv, err := strconv.Atoi(idPaper)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	fmt.Printf("ID Paper %v", idPaper)

	data, err := h.laporanController.GetAllChildFromPaper(idPaperConv)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.ResponseTest{
		Data: data,
	}

	c.JSON(200, resp)
}

func (h *handler) GetAllPaper(c *gin.Context) {
	var request models.PaperFilter
	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	paper, count, err := h.laporanController.GetAllPaper(request)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.ResponseTest{
		Data:  paper,
		Count: count,
	}
	c.JSON(200, resp)
}

func (h *handler) GetAllArtikel(c *gin.Context) {
	var request models.PaperFilter
	err := c.BindJSON(&request)
	if err != nil {
		util.ErrorWrapper(err, 401, c)
		return
	}

	fmt.Printf("LIMIT OFFSET, %v", request.Limit)

	artikel, count, err := h.laporanController.GetAllArtikel(request)
	if err != nil {
		util.ErrorWrapper(err, 500, c)
		return
	}

	resp := models.ResponseTest{
		Data:  artikel,
		Count: count,
	}
	c.JSON(200, resp)
}
