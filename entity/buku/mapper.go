package buku

import (
	"fmt"
	"strings"

	"github.com/ikalkali/rbti-go/entity/models"
)

func mapperQueryBuilder(where string, obj interface{}, tempQuery string, searchQ string) (string, string) {
	if where == "where" {
		trimString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(obj)), ","), "[]")
		tempQuery = tempQuery + fmt.Sprintf(" and %v in (%v)",searchQ, trimString)
	} else {
		where = "where"
		trimString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(obj)), ","), "[]")
		tempQuery = tempQuery + fmt.Sprintf("%v in (%v)", searchQ, trimString)
	}

	return tempQuery, where
}

func getAllJudulQueryBuilder(filter models.FilterBuku) (string) {
	var (
		where string
		tempQuery string
		regexSearch bool
	)

	if len(filter.Kategori) > 0 {
		tempQuery, where = mapperQueryBuilder(where, filter.Kategori, tempQuery, "id_kategori")
	}

	if len(filter.Jenis) > 0 {
		tempQuery, where = mapperQueryBuilder(where, filter.Jenis, tempQuery, "jenis")
	}

	if filter.Judul != "" {
		if where == "where" {
			tempQuery = tempQuery + fmt.Sprintf(" and judul ilike '%%%v%%'", filter.Judul)
		} else {
			where = "where"
			tempQuery = tempQuery + fmt.Sprintf("judul ilike '%%%v%%'", filter.Judul)
		}

		regexSearch = true
	}

	if filter.Penulis != "" {
		if where == "where" {
			if regexSearch {
				tempQuery = tempQuery + fmt.Sprintf(" or penulis ilike '%%%v%%'", filter.Penulis)
			} else {
				tempQuery = tempQuery + fmt.Sprintf(" and penulis ilike '%%%v%%'", filter.Penulis)
			}
		} else {
			where = "where"
			tempQuery = tempQuery + fmt.Sprintf("penulis ilike '%%%v%%'", filter.Penulis)
		}
	}

	if filter.Penerbit != "" {
		if where == "where" {
			if regexSearch {
				tempQuery = tempQuery + fmt.Sprintf(" or penerbit ilike '%%%v%%'", filter.Penerbit)
			} else {
				tempQuery = tempQuery + fmt.Sprintf(" and penerbit ilike '%%%v%%'", filter.Penerbit)
			}
		} else {
			where = "where"
			tempQuery = tempQuery + fmt.Sprintf("penerbit ilike '%%%v%%'", filter.Penerbit)
		}
	}

	return fmt.Sprintf("%v %v", where, tempQuery)
}