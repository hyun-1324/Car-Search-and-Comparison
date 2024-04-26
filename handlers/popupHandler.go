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

// processes a car model name from HTTP requests, retrieves corresponding information for display and cookies, and renders the result page.
func PopupHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form request: ", err)
		http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
		return
	}
	carName := r.FormValue("specifications")

	allCarsInfo, err := apidata.ProcessedApiData()
	if err != nil {
		log.Println("Error fetching data: ", err)
		http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
		return
	}

	aCarInfo := findACarInfoByName(carName, allCarsInfo)

	err = addCookieFromACarInfo(w, r, aCarInfo, allCarsInfo)
	// Set cookie to expire if cookie data can't be processed
	if err != nil {
		log.Println("Error processing cookie data: ", err)
		http.SetCookie(w, &http.Cookie{
			Name: "searchData",
			Expires: time.Date(2019, 6, 5, 11,
				35, 04, 0, time.UTC),
		})
	}

	tmp1 := template.Must(template.ParseFiles("popup.html"))
	tmp1.Execute(w, aCarInfo)

}

func findACarInfoByName(carName string, allCarsInfo []apidata.ProcessedModel) apidata.ProcessedModel {
	var carInfo apidata.ProcessedModel

	for _, model := range allCarsInfo {
		if model.Name == carName {
			carInfo = model
		}
	}

	return carInfo
}

func addCookieFromACarInfo(w http.ResponseWriter, r *http.Request, acarInfo apidata.ProcessedModel, allCarsInfo []apidata.ProcessedModel) error {
	cookie, err := r.Cookie("searchData")
	if err == http.ErrNoCookie {
		makeANewCookie(w, allCarsInfo)
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
	if err := json.Unmarshal([]byte(string(decodedCookieValue)), &data); err != nil {
		return err
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
