package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func comparisonHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	models := r.Form["carmodelName"]

	carInfoForDisplay := findcarinfobyname(models)
	manufacturerInfo := findmanuinfobyname(models)
	categoryInfo := findcategoryinfobyname(models)

	cookie, err := r.Cookie("searchData")
	if err == http.ErrNoCookie {
		emptycookie := CookieData{
			Manufacturer: map[string]int{},
			Category:     map[string]int{},
		}
		mashaledJson, err := json.Marshal(emptycookie)
		if err != nil {
			fmt.Println("Error unmarshaling cookie:", err)
		}

		encodedCookieValue := base64.StdEncoding.EncodeToString([]byte(mashaledJson))

		http.SetCookie(w, &http.Cookie{
			Name:  "searchData",
			Value: encodedCookieValue,
			Path:  "/",
		})

		http.Error(w, "Error parsing data", http.StatusBadRequest)
	}

	decodedcookievalue, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		fmt.Println("Error base64 decoding cookie:", err)
		return
	}

	var data CookieData
	if err := json.Unmarshal([]byte(decodedcookievalue), &data); err != nil {
		fmt.Println("Error unmarshaling cookie:", err)
		return
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
		fmt.Println("Error marshaling cookie:", err)
		return
	}

	encodedJSON := base64.StdEncoding.EncodeToString([]byte(updatedJSON))

	http.SetCookie(w, &http.Cookie{
		Name:  "searchData",
		Value: encodedJSON,
		Path:  "/",
	})

	tmp1 := template.Must(template.ParseFiles("comparison.html"))
	tmp1.Execute(w, carInfoForDisplay)

}

func findcarinfobyname(models []string) []ProcessedModel {
	var carinfo []ProcessedModel
	allCarsinfo, err := processedApiData()
	if err != nil {
		fmt.Println("Error parsing data")
		return nil
	}

	for _, model1 := range models {
		for _, model2 := range allCarsinfo {
			if model1 == model2.Name {
				carinfo = append(carinfo, model2)
			}
		}
	}
	return carinfo
}

func findmanuinfobyname(models []string) []string {
	var manufacturerInfo []string
	allCarsInfo, err := processedApiData()
	if err != nil {
		fmt.Println("Error parsing data")
		return nil
	}

	for _, model1 := range models {
		for _, model2 := range allCarsInfo {
			if model1 == model2.Name {
				manufacturerInfo = append(manufacturerInfo, model2.ManufacturerName)
			}
		}
	}
	return manufacturerInfo
}

func findcategoryinfobyname(models []string) []string {
	var categoryInfo []string
	allCarsInfo, err := processedApiData()
	if err != nil {
		fmt.Println("Error parsing data")
		return nil
	}

	for _, model1 := range models {
		for _, model2 := range allCarsInfo {
			if model1 == model2.Name {
				categoryInfo = append(categoryInfo, model2.CategoryName)
			}
		}
	}
	return categoryInfo
}
