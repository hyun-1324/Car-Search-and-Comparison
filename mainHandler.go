package main

import (
	"encoding/json"
	"html/template"
	"math/rand"
	"net/http"
	"sort"
)

type CookieData struct {
	Manufacturer map[string]int `json:"manufacturer"`
	Categories   map[string]int `json:"category"`
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var bannerData ProcessedModel
	var manufacturerlist []string
	var categorylist []string
	modeldata, err := processedApiData()
	if err != nil {
		http.Error(w, "Error parsing data", http.StatusBadRequest)
	}

	manufacturerlist = findManufacturerlist(modeldata)
	categorylist = findCategorylist(modeldata)

	cookie, err := r.Cookie("searchData")
	if err == http.ErrNoCookie {
		randomNumber := rand.Intn(10)

		bannerData = modeldata[randomNumber]

	} else if err != nil {
		return

	} else {
		var Cookiedata CookieData
		if err := json.Unmarshal([]byte(cookie.Value), &Cookiedata); err != nil {
			http.Error(w, "Error parsing data", http.StatusBadRequest)
		}
		mostSearchedManufacturer, manufacturerSearchnum := maxCookiedataKey(Cookiedata.Manufacturer)
		mostSearchedCategory, categorySearchnum := maxCookiedataKey(Cookiedata.Categories)

		bannerData, err = findMostSearchedData(mostSearchedManufacturer, mostSearchedCategory, manufacturerSearchnum, categorySearchnum)
		if err != nil {
			http.Error(w, "Error parsing data", http.StatusBadRequest)
		}
	}

	tmp1 := template.Must(template.ParseFiles("index.html"))
	tmp1.Execute(w, struct {
		Banner        ProcessedModel
		Manufacturers []string
		Categories    []string
	}{
		Banner:        bannerData,
		Manufacturers: manufacturerlist,
		Categories:    categorylist,
	})

}

func findManufacturerlist(modeldata []ProcessedModel) []string {
	var manufacturerlist []string
	for _, model := range modeldata {
		if !contains(manufacturerlist, model.ManufacturerName) {
			manufacturerlist = append(manufacturerlist, model.ManufacturerName)
		}
	}

	sort.Strings(manufacturerlist)

	return manufacturerlist
}

func findCategorylist(modeldata []ProcessedModel) []string {
	var categorylist []string
	for _, model := range modeldata {
		if !contains(categorylist, model.CategoryName) {
			categorylist = append(categorylist, model.CategoryName)
		}
	}

	sort.Strings(categorylist)

	return categorylist
}

func contains(array []string, str string) bool {
	for _, item := range array {
		if item == str {
			return true
		}
	}
	return false
}

func maxCookiedataKey(data map[string]int) (string, int) {
	var maxKey string
	maxValue := 0
	for key, value := range data {
		if value > maxValue {
			maxValue = value
			maxKey = key
		}
	}
	if maxValue == 0 {
		return "", 0
	}
	return maxKey, maxValue
}

func findMostSearchedData(mostSearchedManufacturer string, mostSearchedCategory string, manufacturerSearchnum int, categorySearchnum int) (result ProcessedModel, err error) {
	var bannerData ProcessedModel
	var categorysave []int
	modeldata, err := processedApiData()
	if err != nil {
		return ProcessedModel{}, err
	}

	if mostSearchedManufacturer != "" && mostSearchedCategory != "" {
		if manufacturerSearchnum >= categorySearchnum {
			for index, model := range modeldata {
				if model.ManufacturerName == mostSearchedManufacturer {
					bannerData = modeldata[index]
					break
				}
			}
		} else {
			for index, model := range modeldata {
				if model.ManufacturerName == mostSearchedManufacturer && model.CategoryName == mostSearchedCategory {
					bannerData = modeldata[index]
					break
				} else if model.CategoryName == mostSearchedCategory {
					bannerData = modeldata[index]
				}
			}
		}

	} else if mostSearchedManufacturer == "" && mostSearchedCategory != "" {
		for _, model := range modeldata {
			if model.CategoryName == mostSearchedCategory {
				categorysave = append(categorysave, model.Id)
			}
		}

		randomIndex := rand.Intn(len(categorysave))
		bannerData = modeldata[categorysave[randomIndex]-1]

	} else if mostSearchedManufacturer != "" && mostSearchedCategory == "" {
		for index, model := range modeldata {
			if model.ManufacturerName == mostSearchedManufacturer {
				bannerData = modeldata[index]
				break
			}
		}
	}

	return bannerData, nil
}
