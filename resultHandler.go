package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func resultHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	selectedManufacturer := r.FormValue("manufacturer")
	selectedCategory := r.FormValue("category")
	carinfoforDisplay := findcarinformation(selectedManufacturer, selectedCategory)

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
	if err := json.Unmarshal([]byte(string(decodedcookievalue)), &data); err != nil {
		fmt.Println("Error unmarshaling cookie:", err)
		return
	}

	_, exists := data.Manufacturer[selectedManufacturer]
	if exists {
		data.Manufacturer[selectedManufacturer]++
	} else {
		if selectedManufacturer != "empty" {
			data.Manufacturer[selectedManufacturer] = 1
		}
	}

	_, exists1 := data.Category[selectedCategory]
	if exists1 {
		data.Category[selectedCategory]++
	} else {
		if selectedCategory != "empty" {
			data.Category[selectedCategory] = 1
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

	tmp1 := template.Must(template.ParseFiles("search_result.html"))
	tmp1.Execute(w, carinfoforDisplay)

}

func findcarinformation(selectedManufacturer string, selectedCategory string) []ProcessedModel {
	var carinfo []ProcessedModel
	allCarsinfo, err := processedApiData()
	if err != nil {
		fmt.Println("Error parsing data")
		return nil
	}

	for _, model := range allCarsinfo {
		if model.ManufacturerName == selectedManufacturer && model.CategoryName == selectedCategory {
			carinfo = append(carinfo, model)
		}
	}

	for _, model := range allCarsinfo {
		if model.ManufacturerName != selectedManufacturer && model.CategoryName == selectedCategory || model.ManufacturerName == selectedManufacturer && model.CategoryName != selectedCategory {
			carinfo = append(carinfo, model)
		}
	}
	return carinfo
}
