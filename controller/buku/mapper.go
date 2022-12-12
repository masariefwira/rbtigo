package buku

import (
	"fmt"

	"github.com/ikalkali/rbti-go/entity/models"
)

func mapperJudulToJudulElastic(b *bukuController, input *models.InputJudulBaruRequest, idJudul int) (models.JudulElastic, error) {
	var (
		bukuElastic models.JudulElastic
	)

	bukuElastic.Id = idJudul
	bukuElastic.Judul = input.DetailJudul.Judul
	bukuElastic.Bahasa = input.DetailJudul.Bahasa
	bukuElastic.Jenis = input.DetailJudul.Jenis
	bukuElastic.Tahun = input.DetailJudul.Tahun
	bukuElastic.Penulis = input.DetailJudul.Penulis
	bukuElastic.Penerbit = input.DetailJudul.Penerbit
	bukuElastic.IDKategori = input.DetailJudul.IDKategori
	bukuElastic.Tipe = "buku"
	bukuElastic.JumlahTersedia = input.Jumlah
	bukuElastic.JumlahTotal = input.Jumlah

	kategori, err := b.bukuDb.GetKategoriBuku(input.DetailJudul.IDKategori)
	if err != nil {
		return bukuElastic, err
	}

	bukuElastic.Kategori = kategori

	return bukuElastic, nil
}

func mapperUpdateToElastic(b *bukuController, input *models.Judul) (models.JudulElastic, error) {
	var (
		bukuElastic models.JudulElastic
	)

	bukuElastic.Id = input.Id
	bukuElastic.Judul = input.Judul
	bukuElastic.Bahasa = input.Bahasa
	bukuElastic.Jenis = input.Jenis
	bukuElastic.Tahun = input.Tahun
	bukuElastic.Penulis = input.Penulis
	bukuElastic.Penerbit = input.Penerbit
	bukuElastic.IDKategori = input.IDKategori
	bukuElastic.Tipe = "buku"

	oldIndex, err := b.bukuEs.GetJudulByID(models.ElasticFilter{ID: input.Id, Jenis: "buku"})
	if err != nil {
		return bukuElastic, err
	}
	fmt.Printf("OLD INDEX %+v", oldIndex)

	bukuElastic.JumlahTersedia = oldIndex.JumlahTersedia
	bukuElastic.JumlahTotal = oldIndex.JumlahTotal
	

	kategori, err := b.bukuDb.GetKategoriBuku(input.IDKategori)
	if err != nil {
		return bukuElastic, err
	}

	bukuElastic.Kategori = kategori

	return bukuElastic, nil
}

func MapAndValidateUpdateJudul(input *models.Judul, oldJudul *models.Judul) {
	if input.Judul == "" {
		input.Judul = oldJudul.Judul
	}

	if input.Bahasa == "" {
		input.Bahasa = oldJudul.Bahasa
	}

	if input.Tahun == 0 {
		input.Tahun = oldJudul.Tahun
	}

	if input.Penerbit == "" {
		input.Penerbit = oldJudul.Penerbit
	}

	if input.Filename == nil {
		input.Filename = oldJudul.Filename
	}

	if input.Foto == nil {
		input.Foto = oldJudul.Foto
	}

	if input.Jenis == 0 {
		input.Jenis = oldJudul.Jenis
	}

	if input.IDKategori == 0 {
		input.IDKategori = oldJudul.IDKategori
	}
}