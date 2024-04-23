package handlers

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	apidata "cars/apiDataProcess"
)

// processes car model names from HTTP requests, retrieves corresponding information for display and cookies, and renders the comparison page.
func ComparisonHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	models := r.Form["carmodelName"]

	allCarsInfo, err := apidata.ProcessedApiData()
	if err != nil {
		log.Println("Error fetching data: ", err)
		http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
		return
	}

	carInfoForDisplay := findCarsInfoByName(models, allCarsInfo)
	manufacturerInfo := findManuInfoByName(models, allCarsInfo)
	categoryInfo := findCategoryInfoByName(models, allCarsInfo)

	err = addCookieFromManufacturersAndCategoriesInfo(w, r, manufacturerInfo, categoryInfo, allCarsInfo)
	// Set cookie to expire if cookie data can't be processed
	if err != nil {
		log.Println("Error processing cookie data: ", err)
		http.SetCookie(w, &http.Cookie{
			Name: "searchData",
			Expires: time.Date(2019, 6, 5, 11,
				35, 04, 0, time.UTC),
		})
	}

	tmp1 := template.Must(template.ParseFiles("comparison.html"))
	tmp1.Execute(w, carInfoForDisplay)

}

func findCarsInfoByName(models []string, allCarsInfo []apidata.ProcessedModel) []apidata.ProcessedModel {
	var carInfo []apidata.ProcessedModel

	for _, model1 := range models {
		for _, model2 := range allCarsInfo {
			if model1 == model2.Name {
				carInfo = append(carInfo, model2)
			}
		}
	}
	return carInfo
}

func findManuInfoByName(models []string, allCarsInfo []apidata.ProcessedModel) []string {
	var manufacturerInfo []string

	for _, model1 := range models {
		for _, model2 := range allCarsInfo {
			if model1 == model2.Name {
				manufacturerInfo = append(manufacturerInfo, model2.ManufacturerName)
			}
		}
	}
	return manufacturerInfo
}

func findCategoryInfoByName(models []string, allCarsInfo []apidata.ProcessedModel) []string {
	var categoryInfo []string

	for _, model1 := range models {
		for _, model2 := range allCarsInfo {
			if model1 == model2.Name {
				categoryInfo = append(categoryInfo, model2.CategoryName)
			}
		}
	}
	return categoryInfo
}

func addCookieFromManufacturersAndCategoriesInfo(w http.ResponseWriter, r *http.Request, manufacturerInfo []string, categoryInfo []string, allCarsInfo []apidata.ProcessedModel) error {
	cookie, err := r.Cookie("searchData")
	if err == http.ErrNoCookie {
		err = makeANewCookie(w, allCarsInfo)
		if err != nil {
			return err
		}
		cookie, _ = r.Cookie("searchData")
	}

	decodedCookieValue, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return err
	}

	var data CookieData
	if err := json.Unmarshal([]byte(decodedCookieValue), &data); err != nil {
		return err
	}

	for _, manufacturer := range manufacturerInfo {
		_, exists := data.Manufacturer[manufacturer]
		if exists {
			data.Manufacturer[manufacturer]++
		} else {
			data.Manufacturer[manufacturer] = 1
		}
	}

	for _, category := range categoryInfo {
		_, exists1 := data.Category[category]
		if exists1 {
			data.Category[category]++
		} else {
			data.Category[category] = 1
		}
	}

	updatedJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	encodedJSON := base64.StdEncoding.EncodeToString([]byte(updatedJSON))

	http.SetCookie(w, &http.Cookie{
		Name:  "searchData",
		Value: encodedJSON,
		Path:  "/",
	})

	return nil
}
