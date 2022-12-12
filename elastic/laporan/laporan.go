package laporan

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/ikalkali/rbti-go/entity/models"
	"github.com/olivere/elastic/v7"
)

type LaporanElasticInterface interface {
	InsertLaporan(input models.JudulElastic) (error)
	SearchLaporan(input models.ElasticFilter) ([]models.JudulElastic, int, error)
	InsertKaryaTulis(input models.JudulElastic) (error)
	GetKaryaTulisByID(id int) (models.JudulElastic, error)
	SearchKaryaTulis(input models.ElasticFilter) ([]models.JudulElastic, int, error)
}

type laporan struct {
	client *elastic.Client
}

func NewElasticRepo(client *elastic.Client) *laporan {
	return &laporan{client: client}
}

func (l *laporan) InsertLaporan(input models.JudulElastic) (error) {
	idJudul := strconv.Itoa(input.Id)

	resp, err := l.client.Index().Index("laporan").Id(idJudul).BodyJson(&input).Do(context.Background())
	if err != nil {
		return err
	}

	log.Print(resp)
	return nil
}

func (l *laporan) SearchLaporan(input models.ElasticFilter) ([]models.JudulElastic, int, error) {
	var (
		shouldQueries []elastic.Query
		mustQueries []elastic.Query
		resp []models.JudulElastic
	)

	

	if input.Jenis != "" {
		mustQueries = append(mustQueries, elastic.NewMatchQuery("tipe", input.Jenis))
	}

	if len(input.IDKategori) > 0 {
		var temp []interface{}
		for _, idKategori := range input.IDKategori {
			temp = append(temp, idKategori)
		}
		mustQueries = append(mustQueries, elastic.NewTermsQuery("id_kategori", temp...))
	}
	

	shouldQueries = append(shouldQueries, elastic.NewPrefixQuery("judul", input.Query).CaseInsensitive(true))
	shouldQueries = append(shouldQueries, elastic.NewPrefixQuery("penulis", input.Query).CaseInsensitive(true))
	shouldQueries = append(shouldQueries, elastic.NewPrefixQuery("kategori", input.Query).CaseInsensitive(true))
	shouldQueries = append(shouldQueries, elastic.NewPrefixQuery("nim", input.Query).CaseInsensitive(true))

	shouldQueries = append(shouldQueries, elastic.NewMatchQuery("judul", input.Query))
	shouldQueries = append(shouldQueries, elastic.NewMatchQuery("penulis", input.Query))
	shouldQueries = append(shouldQueries, elastic.NewMatchQuery("kategori", input.Query))

	q := elastic.NewBoolQuery().Must(mustQueries...).Should(shouldQueries...).MinimumNumberShouldMatch(1)

	src, err := q.Source()
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	search := l.client.Search().Index("laporan").Query(q)
	fmt.Printf("SIZE : %v FROM : %v", input.Size, input.From)
	search.Size(input.Size)
	search.From(input.From)


	result, err := search.TrackTotalHits(true).Do(context.Background())
	if err != nil {
		return resp, 0,  err
	}

	searchCount := result.TotalHits()


	for _, hit := range result.Hits.Hits {
		temp := models.JudulElastic{}
		
		err := json.Unmarshal(hit.Source, &temp)
		if err != nil {
			return resp,0, err
		}

		resp = append(resp, temp)
	}

	return resp, int(searchCount), nil
}

func (l *laporan) SearchKaryaTulis(input models.ElasticFilter) ([]models.JudulElastic, int, error) {
	var (
		shouldQueries []elastic.Query
		mustQueries []elastic.Query
		resp []models.JudulElastic
	)


	if len(input.IDKategori) > 0 {
		var temp []interface{}
		for _, idKategori := range input.IDKategori {
			temp = append(temp, idKategori)
		}
		mustQueries = append(mustQueries, elastic.NewTermsQuery("id_kategori", temp...))
	}
	

	shouldQueries = append(shouldQueries, elastic.NewPrefixQuery("judul", input.Query).CaseInsensitive(true))
	shouldQueries = append(shouldQueries, elastic.NewPrefixQuery("penulis", input.Query).CaseInsensitive(true))
	shouldQueries = append(shouldQueries, elastic.NewPrefixQuery("kategori", input.Query).CaseInsensitive(true))

	shouldQueries = append(shouldQueries, elastic.NewMatchQuery("judul", input.Query))
	shouldQueries = append(shouldQueries, elastic.NewMatchQuery("penulis", input.Query))
	shouldQueries = append(shouldQueries, elastic.NewMatchQuery("kategori", input.Query))

	q := elastic.NewBoolQuery().Must(mustQueries...).Should(shouldQueries...).MinimumNumberShouldMatch(1)

	src, err := q.Source()
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	search := l.client.Search().Index("karya_tulis").Query(q)
	fmt.Printf("SIZE : %v FROM : %v", input.Size, input.From)
	search.Size(input.Size)
	search.From(input.From)


	result, err := search.TrackTotalHits(true).Do(context.Background())
	if err != nil {
		return resp, 0,  err
	}

	searchCount := result.TotalHits()


	for _, hit := range result.Hits.Hits {
		temp := models.JudulElastic{}
		
		err := json.Unmarshal(hit.Source, &temp)
		if err != nil {
			return resp,0, err
		}

		resp = append(resp, temp)
	}

	return resp, int(searchCount), nil
}

func (l *laporan) GetKaryaTulisByID(id int) (models.JudulElastic, error) {
	var (
		mustQueries []elastic.Query
		resp models.JudulElastic
		err error
		search *elastic.SearchService
	)

	mustQueries = append(mustQueries, elastic.NewMatchQuery("id", id))

	q := elastic.NewBoolQuery().Must(mustQueries...)


	search = l.client.Search().Index("karya_tulis").Query(q)

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

func (l *laporan) InsertKaryaTulis(input models.JudulElastic) (error) {
	id := strconv.Itoa(input.Id)

	resp, err := l.client.Index().Index("karya_tulis").Id(id).BodyJson(&input).Do(context.Background())
	if err != nil {
		return err
	}

	log.Print(resp)
	return nil
}

