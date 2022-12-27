package buku

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ikalkali/rbti-go/entity/models"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

type BukuElasticInterface interface {
	InsertBuku(input models.JudulElastic)
	SearchBuku(input models.ElasticFilter) ([]models.JudulElastic, int, error)
	ChangeBukuStock(idJudul int, action string) error
	GetJudulByID(input models.ElasticFilter) (models.JudulElastic, error)
	DeleteBuku(input models.ElasticFilter) error
}

type buku struct {
	client *elastic.Client
}

func NewElasticRepo(client *elastic.Client) *buku {
	return &buku{client: client}
}

func (b *buku) InsertBuku(input models.JudulElastic) {
	idJudul := strconv.Itoa(input.Id)

	resp, err := b.client.Index().Index("buku").Id(idJudul).BodyJson(&input).Do(context.Background())
	if err != nil {
		log.Error(err)
	}

	fmt.Printf("ERROR ELASTIC :%v\n", err)

	log.Print(resp)
}

func (b *buku) ChangeBukuStock(idJudul int, action string) error {
	var (
		resp models.JudulElastic
	)
	idJudulConv := strconv.Itoa(idJudul)

	// get past index
	pastIndex, err := b.client.Get().Index("buku").Id(idJudulConv).Do(context.Background())
	if err != nil {
		return err
	}

	if pastIndex.Found {
		json.Unmarshal(pastIndex.Source, &resp)
	} else {
		return errors.New("[ChangeBukuStock] no index found with the given judul ID")
	}

	switch action {
	case "increment":
		fmt.Println("INCREMENT CALLED")
		resp.JumlahTersedia = resp.JumlahTersedia + 1
	case "decrement":
		fmt.Println("DECREMENT CALLED")
		resp.JumlahTersedia = resp.JumlahTersedia - 1
	default:
		return errors.New("[ChangeBukuStock] invalid action type, the available action types are 'increment' and 'decrement'")
	}

	_, err = b.client.Index().Index("buku").Id(idJudulConv).BodyJson(&resp).Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (b *buku) SearchBuku(input models.ElasticFilter) ([]models.JudulElastic, int, error) {
	var (
		shouldQueries []elastic.Query
		mustQueries   []elastic.Query
		resp          []models.JudulElastic
	)

	if len(input.IDKategori) > 0 {
		var temp []interface{}
		for _, idKategori := range input.IDKategori {
			temp = append(temp, idKategori)
		}
		mustQueries = append(mustQueries, elastic.NewTermsQuery("id_kategori", temp...))
	}

	if input.JenisPinjam != 0 {
		mustQueries = append(mustQueries, elastic.NewMatchQuery("jenis", input.JenisPinjam))
		if input.BisaDipinjam {
			mustQueries = append(mustQueries, elastic.NewRangeQuery("jumlah_tersedia").Gt(0))
		}
	}

	shouldQueries = append(shouldQueries, elastic.NewPrefixQuery("judul", input.Query).CaseInsensitive(true))
	shouldQueries = append(shouldQueries, elastic.NewPrefixQuery("penerbit", input.Query).CaseInsensitive(true))
	shouldQueries = append(shouldQueries, elastic.NewPrefixQuery("penulis", input.Query).CaseInsensitive(true))
	shouldQueries = append(shouldQueries, elastic.NewPrefixQuery("kategori", input.Query).CaseInsensitive(true))

	shouldQueries = append(shouldQueries, elastic.NewMatchQuery("judul", input.Query))
	shouldQueries = append(shouldQueries, elastic.NewMatchQuery("penerbit", input.Query))
	shouldQueries = append(shouldQueries, elastic.NewMatchQuery("penulis", input.Query))
	shouldQueries = append(shouldQueries, elastic.NewMatchQuery("kategori", input.Query))

	q := elastic.NewBoolQuery().Should(shouldQueries...).Must(mustQueries...).MinimumNumberShouldMatch(1)

	// src, err := q.Source()
	// if err != nil {
	// 	panic(err)
	// }
	// data, err := json.MarshalIndent(src, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(data))

	search := b.client.Search().Index("buku").Query(q)
	fmt.Printf("SIZE : %v FROM : %v", input.Size, input.From)
	search.Size(input.Size)
	search.From(input.From)

	result, err := search.TrackTotalHits(true).Do(context.Background())
	if err != nil {
		return resp, 0, err
	}

	searchCount := result.TotalHits()

	fmt.Printf("TOTAL HITS BUKu %+v\n", result.TotalHits())

	for _, hit := range result.Hits.Hits {
		temp := models.JudulElastic{}

		err := json.Unmarshal(hit.Source, &temp)
		if err != nil {
			return resp, int(searchCount), err
		}

		resp = append(resp, temp)
	}

	return resp, int(searchCount), nil
}

func (b *buku) GetJudulByID(input models.ElasticFilter) (models.JudulElastic, error) {
	var (
		mustQueries []elastic.Query
		resp        models.JudulElastic
		err         error
		search      *elastic.SearchService
	)

	mustQueries = append(mustQueries, elastic.NewMatchQuery("id", input.ID))
	mustQueries = append(mustQueries, elastic.NewMatchQuery("tipe", input.Jenis))

	q := elastic.NewBoolQuery().Must(mustQueries...)

	if strings.ToLower(input.Jenis) == "buku" {
		search = b.client.Search().Index("buku").Query(q)
	} else {
		search = b.client.Search().Index("laporan").Query(q)
	}

	result, err := search.Do(context.Background())
	if err != nil {
		return resp, err
	}

	for _, hit := range result.Hits.Hits {
		temp := models.JudulElastic{}

		err := json.Unmarshal(hit.Source, &temp)
		if err != nil {
			return resp, err
		}

		resp = temp
	}

	return resp, nil
}

func (b *buku) DeleteBuku(input models.ElasticFilter) error {
	indexName := strings.ToLower(input.Jenis)
	indexElastic := ""
	if indexName == "buku" {
		indexElastic = "buku"
	} else {
		indexElastic = "laporan"
	}

	idStr := fmt.Sprint(input.ID)

	_, err := b.client.Delete().Index(indexElastic).Id(idStr).Refresh("true").Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}
