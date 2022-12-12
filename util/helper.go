package util

import (
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ikalkali/rbti-go/entity/models"
)

func HitungDenda(tanggalPinjam time.Time) int {
	now := time.Now()
	timeSincePinjam := now.Sub(tanggalPinjam).Hours()

	if timeSincePinjam < DURASI_PEMINJAMAN.Hours() {
		return 0
	}

	totalKeterlambatan := math.Round((timeSincePinjam - DURASI_PEMINJAMAN.Hours())/24)
	denda := int(totalKeterlambatan) * DENDA_PER_HARI
	return denda
}

func HitungDendaByTenggat(tenggat time.Time) int {
	now := time.Now()
	timeSinceTenggat := now.Sub(tenggat).Hours()

	if timeSinceTenggat < 0 {
		return 0
	}

	return int(timeSinceTenggat) * DENDA_PER_HARI
}

func ErrorWrapper(err error, statusCode int, c *gin.Context) () {
	var (
		resp models.Response
	)

	resp.Errors = append(resp.Errors, err.Error())

	c.JSON(statusCode, resp)
	return
}