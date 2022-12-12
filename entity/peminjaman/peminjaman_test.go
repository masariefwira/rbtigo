package peminjaman

import (
	"testing"

	models "github.com/ikalkali/rbti-go/entity/models"
	"gorm.io/gorm"
)

func Test_peminjaman_GetJudulDetailBukuDipinjamByIdPeminjaman(t *testing.T) {
	type args struct {
		input models.DetailBukuPeminjamanFilter
	}
	tests := []struct {
		name    string
		p       *peminjaman
		args    args
		want    []models.DetailPeminjaman
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			p: &peminjaman{
				db: &gorm.DB{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// _, err := tt.p.GetJudulDetailBukuDipinjamByIdPeminjaman(tt.args.input)
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("peminjaman.GetJudulDetailBukuDipinjamByIdPeminjaman() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
		})
	}
}
