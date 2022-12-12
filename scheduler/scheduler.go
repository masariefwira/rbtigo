package scheduler

import (
	"fmt"
	"time"

	"github.com/ikalkali/rbti-go/controller/peminjaman"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SchedulerInterface interface{}

type scheduler struct {
	cron        *cron.Cron
	db          *gorm.DB
	cPeminjaman peminjaman.PeminjamanControllerInterface
}

func NewScheduler(cron *cron.Cron, db *gorm.DB, cPeminjaman peminjaman.PeminjamanControllerInterface) *scheduler {
	return &scheduler{cron, db, cPeminjaman}
}

func (s *scheduler) LogEveryOneHour() {
	query := "INSERT INTO log (message) VALUES (?)"

	timeString := fmt.Sprint(time.Now().Hour())

	timePrint := fmt.Sprintf("Jam %v", timeString)

	err := s.db.Exec(query, timePrint).Error
	if err != nil {
		logrus.Error(err)
		return
	}
}
