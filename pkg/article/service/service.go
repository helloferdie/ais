package service

import (
	"ais/lib/libresponse"
	"ais/lib/libvalidator"
	"ais/pkg/article"
	"ais/pkg/article/request"
	"ais/pkg/article/response"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/helloferdie/golib/libdb"
	"github.com/helloferdie/golib/libredis"
	"github.com/jinzhu/copier"
)

// List - List articles
func List(r *request.List, format map[string]interface{}) *libresponse.Response {
	resp, err := libvalidator.Validate(r)
	if err != nil {
		return resp
	}

	d, _ := libdb.Open("")
	defer d.Close()

	// Get list ID of relevant articles
	condition := "AND deleted_at IS NULL "
	conditionVal := map[string]interface{}{}

	// Filter author with submatch string
	if r.Author != "" {
		condition += "AND author LIKE :author "
		conditionVal["author"] = "%" + r.Author + "%"
	}

	// Filter keyword on title and body with minimum chars (default database configuration 3 chars)
	if r.Query != "" {
		condition += "AND MATCH(title, body) AGAINST (:keyword) > 0 "
		conditionVal["keyword"] = r.Query
	}

	listID := []libdb.ModelID{}
	totalItems, err := libdb.ListByField(d, &listID, conditionVal, condition, "article", "id", "id", &libdb.ModelPaginationRequest{
		ShowAll:          false,
		Page:             r.Page,
		ItemsPerPage:     r.ItemsPerPage,
		OrderByField:     "created_at",
		OrderByDirection: "desc",
	})
	if err != nil {
		return resp.ErrorList()
	}

	result := make([]interface{}, len(listID)) // Allocate

	// When relevant articles not found, skip following
	if totalItems > 0 {
		// Generate sequence map from list ID
		listSequence := map[int64]int{}
		for k, v := range listID {
			listSequence[v.ID] = k
		}

		// Iterate and get data from cache
		cacheEnable := true
		rd := new(libredis.Client)

		listDBID := []int64{} // Store ID to be query on DB later
		for id, seq := range listSequence {
			// When cache is disable
			if !cacheEnable || (rd.HasInitialize && !rd.Enable) {
				listDBID = append(listDBID, id)
				continue
			}

			tmp := new(article.Model)
			cacheKey := "article_" + strconv.FormatInt(id, 10)
			cacheExist, cacheErr := rd.GetUnmarshal(cacheKey, tmp)
			if cacheErr != nil || !cacheExist {
				if cacheErr != nil {
					// When something error on get cache, disable cache
					cacheEnable = false
				}

				// When error on getting cache or cache not exist
				listDBID = append(listDBID, id)
				continue
			}

			// Assign result using cache data
			result[seq] = response.Article(tmp, format)
		}

		if len(listDBID) > 0 {
			// Iterate listDBID and get data from database
			namedParamID := []string{}
			conditionVal = map[string]interface{}{}
			for k, id := range listDBID {
				// Generate query for select IN condition
				param := "id_" + strconv.Itoa(k)
				conditionVal[param] = id
				namedParamID = append(namedParamID, ":"+param)
			}

			// Query list articles from database
			listTmp := []article.Model{}
			_, err = libdb.List(d, article.TConfig, &listTmp, conditionVal, "AND id IN ("+strings.Join(namedParamID, ", ")+") ", &libdb.ModelPaginationRequest{
				ShowAll: true,
			})
			if err != nil {
				return resp.ErrorList()
			}

			for _, at := range listTmp {
				// Assign query result using db data
				seq := listSequence[at.ID]
				result[seq] = response.Article(&at, format)

				// Save query result to cache
				if cacheEnable {
					cacheKey := "article_" + strconv.FormatInt(at.ID, 10)
					b, _ := json.Marshal(at)
					err = rd.Set(cacheKey, string(b), true)
					if err != nil {
						cacheEnable = false
					}
				}
			}
		}
	}

	resp.Data = map[string]interface{}{
		"items":       result,
		"total_items": totalItems,
		"total_pages": libresponse.TotalPages(r.ItemsPerPage, totalItems),
	}
	return resp.SuccessList()
}

// View - View article
func View(r *request.View, format map[string]interface{}) *libresponse.Response {
	resp, err := libvalidator.Validate(r)
	if err != nil {
		return resp
	}

	d, _ := libdb.Open("")
	defer d.Close()

	at := new(article.Model)

	rd := new(libredis.Client)
	cacheKey := "article_" + strconv.FormatInt(r.ID, 10)
	cacheExist, cacheErr := rd.GetUnmarshal(cacheKey, at)
	if cacheErr != nil || !cacheExist {
		// Cache not exist, get from database instead
		exist, err := libdb.GetByID(d, article.TConfig, at, r.ID)
		if err != nil || !exist {
			return resp.ErrorDataNotFound()
		}

		go RefreshCache(at, true)
	}

	resp.Data = response.Article(at, format)
	return resp.SuccessDefault()
}

// Post - Post new article
func Post(r *request.Post, format map[string]interface{}) *libresponse.Response {
	resp, err := libvalidator.Validate(r)
	if err != nil {
		return resp
	}

	d, _ := libdb.Open("")
	defer d.Close()

	at := new(article.Model)
	at.Author = r.Author
	at.Title = r.Title
	at.Body = r.Body

	err = libdb.Create(d, article.TConfig, at, article.DBMode, true)
	if err != nil {
		return resp.ErrorCreate()
	}

	resp.Data = response.Article(at, format)
	return resp.SuccessCreate()
}

// Update - Update article
func Update(r *request.Update, format map[string]interface{}) *libresponse.Response {
	resp, err := libvalidator.Validate(r)
	if err != nil {
		return resp
	}

	d, _ := libdb.Open("")
	defer d.Close()

	at := new(article.Model)
	exist, err := libdb.GetByID(d, article.TConfig, at, r.ID)
	if err != nil || !exist {
		return resp.ErrorDataNotFound()
	}

	oldArticle := new(article.Model)
	copier.Copy(&oldArticle, &at)
	at.Author = r.Author
	at.Title = r.Title
	at.Body = r.Body
	_, err = libdb.Update(d, article.TConfig, oldArticle, at, article.DBMode, at.ID, true)
	if err != nil {
		return resp.ErrorUpdate()
	}

	go RefreshCache(at, true)

	resp.Data = response.Article(at, format)
	return resp.SuccessUpdate()
}

// Delete - Delete article
func Delete(r *request.Delete, format map[string]interface{}) *libresponse.Response {
	resp, err := libvalidator.Validate(r)
	if err != nil {
		return resp
	}

	d, _ := libdb.Open("")
	defer d.Close()

	at := new(article.Model)
	exist, err := libdb.GetByID(d, article.TConfig, at, r.ID)
	if err != nil || !exist {
		return resp.ErrorDataNotFound()
	}

	err = libdb.Delete(d, article.TConfig, at.ID)
	if err != nil {
		return resp.ErrorDelete()
	}

	go RefreshCache(at, false)
	return resp.SuccessDelete()
}

// RefreshCache - Refresh article data on redis
func RefreshCache(at *article.Model, update bool) {
	rd := new(libredis.Client)
	cacheKey := "article_" + strconv.FormatInt(at.ID, 10)
	if update {
		b, _ := json.Marshal(at)
		rd.Set(cacheKey, string(b), true)
	} else {
		rd.Delete(cacheKey)
	}
}
