package handlers

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"sort"
	"time"

	apidata "cars/apiDataProcess"
)

type CookieData struct {
	Manufacturer map[string]int `json:"manufacturer"`
	Category     map[string]int `json:"category"`
}

// processes data from the API server, retrieving relevant manufacturer and category information for display and cookies, and renders the results on the page.
func MainHandler(w http.ResponseWriter, r *http.Request) {
	var bannerData apidata.ProcessedModel
	var manufacturerList []string
	var categoryList []string

	modelData, err := apidata.ProcessedApiData()
	if err != nil {
		log.Println("Error fetching data: ", err)
		http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
		return
	}

	manufacturerList = findManufacturerList(modelData)
	categoryList = findCategoryList(modelData)

	cookie, err := r.Cookie("searchData")
	if err == http.ErrNoCookie {
		err := makeANewCookie(w, modelData)
		if err != nil {
			log.Println("Error making a new cookie: ", err)
			http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
			return
		}
		randomNumber := rand.Intn(10)
		bannerData = modelData[randomNumber]
	} else {
		bannerData, err = getBannerDataFromCookie(cookie, modelData)
		// Set cookie to expire if error is occured during processing cookie data
		if err != nil {
			log.Println("Error parsing data: ", err)
			http.SetCookie(w, &http.Cookie{
				Name: "searchData",
				Expires: time.Date(2019, 6, 5, 11,
					35, 04, 0, time.UTC),
			})
			err := makeANewCookie(w, modelData)
			if err != nil {
				log.Println("Error making a new cookie: ", err)
				http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
				return
			}
		}
	}

	tmp1 := template.Must(template.ParseFiles("index.html"))
	tmp1.Execute(w, struct {
		Banner        apidata.ProcessedModel
		Manufacturers []string
		Categories    []string
	}{
		Banner:        bannerData,
		Manufacturers: manufacturerList,
		Categories:    categoryList,
	})

}

func findManufacturerList(modelData []apidata.ProcessedModel) []string {
	var manufacturerList []string
	for _, model := range modelData {
		if !slices.Contains(manufacturerList, model.ManufacturerName) {
			manufacturerList = append(manufacturerList, model.ManufacturerName)
		}
	}

	sort.Strings(manufacturerList)

	return manufacturerList
}

func findCategoryList(modelData []apidata.ProcessedModel) []string {
	var categoryList []string
	for _, model := range modelData {
		if !slices.Contains(categoryList, model.CategoryName) {
			categoryList = append(categoryList, model.CategoryName)
		}
	}

	sort.Strings(categoryList)

	return categoryList
}

func maxCookieDataKey(data map[string]int) (string, int) {
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

func findMostSearchedData(mostSearchedManufacturer string, mostSearchedCategory string, manufacturerSearchNum int, categorySearchNum int, modelData []apidata.ProcessedModel) apidata.ProcessedModel {
	var bannerData apidata.ProcessedModel
	var categorySave []int

	if mostSearchedManufacturer != "" && mostSearchedCategory != "" {
		if manufacturerSearchNum >= categorySearchNum {
			for index, model := range modelData {
				if model.ManufacturerName == mostSearchedManufacturer {
					bannerData = modelData[index]
					break
				}
			}
		} else {
			for index, model := range modelData {
				if model.ManufacturerName == mostSearchedManufacturer && model.CategoryName == mostSearchedCategory {
					bannerData = modelData[index]
					break
				} else if model.CategoryName == mostSearchedCategory {
					bannerData = modelData[index]
				}
			}
		}

	} else if mostSearchedManufacturer == "" && mostSearchedCategory != "" {
		for _, model := range modelData {
			if model.CategoryName == mostSearchedCategory {
				categorySave = append(categorySave, model.Id)
			}
		}

		randomIndex := rand.Intn(len(categorySave))
		bannerData = modelData[categorySave[randomIndex]-1]

	} else if mostSearchedManufacturer != "" && mostSearchedCategory == "" {
		for index, model := range modelData {
			if model.ManufacturerName == mostSearchedManufacturer {
				bannerData = modelData[index]
				break
			}
		}
	}

	return bannerData
}

func makeANewCookie(w http.ResponseWriter, modelData []apidata.ProcessedModel) error {
	var randomManufacturerName string
	var randomCategoryName string
	randomNumber := rand.Intn(10) + 1
	for _, model := range modelData {
		if model.Id == randomNumber {
			randomManufacturerName = model.ManufacturerName
			randomCategoryName = model.CategoryName
		}
	}

	emptyCookie := CookieData{
		Manufacturer: map[string]int{randomManufacturerName: 1},
		Category:     map[string]int{randomCategoryName: 1},
	}
	mashaledJson, err := json.Marshal(emptyCookie)
	if err != nil {
		return err
	}

	encodedCookieValue := base64.StdEncoding.EncodeToString([]byte(mashaledJson))

	http.SetCookie(w, &http.Cookie{
		Name:  "searchData",
		Value: encodedCookieValue,
		Path:  "/",
	})

	return nil
}

func getBannerDataFromCookie(cookie *http.Cookie, modelData []apidata.ProcessedModel) (apidata.ProcessedModel, error) {
	var bannerData apidata.ProcessedModel

	decodedCookieValue, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return apidata.ProcessedModel{}, err
	}

	var CookieData CookieData
	if err := json.Unmarshal([]byte(string(decodedCookieValue)), &CookieData); err != nil {
		return apidata.ProcessedModel{}, err
	}
	mostSearchedManufacturer, manufacturerSearchNum := maxCookieDataKey(CookieData.Manufacturer)
	mostSearchedCategory, categorySearchNum := maxCookieDataKey(CookieData.Category)

	bannerData = findMostSearchedData(mostSearchedManufacturer, mostSearchedCategory, manufacturerSearchNum, categorySearchNum, modelData)

	return bannerData, nil

}
