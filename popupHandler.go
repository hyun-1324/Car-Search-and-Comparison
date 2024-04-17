package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func popupHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	carName := r.FormValue("specifications")
	acarInfo := findacarinfobyname(carName)

	addCookieFromAcarinfo(w, r, acarInfo)

	tmp1 := template.Must(template.ParseFiles("popup.html"))
	tmp1.Execute(w, acarInfo)

}

func findacarinfobyname(carName string) ProcessedModel {
	var carinfo ProcessedModel
	allCarsinfo, err := processedApiData()
	if err != nil {
		fmt.Println("Error parsing data")
		return ProcessedModel{}
	}

	for _, model := range allCarsinfo {
		if model.Name == carName {
			carinfo = model
		}
	}

	return carinfo
}

func addCookieFromAcarinfo(w http.ResponseWriter, r *http.Request, acarInfo ProcessedModel) {
	cookie, err := r.Cookie("searchData")
	if err == http.ErrNoCookie {
		makeanewcookie(w)

		cookie, _ = r.Cookie("searchData")

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

	_, exists := data.Manufacturer[acarInfo.ManufacturerName]
	if exists {
		data.Manufacturer[acarInfo.ManufacturerName]++
	} else {
		if acarInfo.ManufacturerName != "empty" {
			data.Manufacturer[acarInfo.ManufacturerName] = 1
		}
	}

	_, exists1 := data.Category[acarInfo.CategoryName]
	if exists1 {
		data.Category[acarInfo.CategoryName]++
	} else {
		if acarInfo.CategoryName != "empty" {
			data.Category[acarInfo.CategoryName] = 1
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
}
