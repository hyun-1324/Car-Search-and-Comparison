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

// processes manufacturers and categories from HTTP requests, retrieves corresponding information for display and cookies, and renders the result page.
func ResultHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form request: ", err)
		http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
		return
	}
	selectedManufacturer := r.FormValue("manufacturer")
	selectedCategory := r.FormValue("category")

	allCarsInfo, err := apidata.ProcessedApiData()
	if err != nil {
		log.Println("Error fetching data: ", err)
		http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
		return
	}

	carsInfoForDisplay := findCarsInfo(selectedManufacturer, selectedCategory, allCarsInfo)

	err = addCookieFromAManufacturerAndCategoryinfo(w, r, selectedManufacturer, selectedCategory, allCarsInfo)
	// Set cookie to expire if cookie data can't be processed
	if err != nil {
		log.Println("Error parsing data: ", err)
		http.SetCookie(w, &http.Cookie{
			Name: "searchData",
			Expires: time.Date(2019, 6, 5, 11,
				35, 04, 0, time.UTC),
		})
	}

	tmp1 := template.Must(template.ParseFiles("search_result.html"))
	tmp1.Execute(w, carsInfoForDisplay)

}

func findCarsInfo(selectedManufacturer string, selectedCategory string, allCarsInfo []apidata.ProcessedModel) []apidata.ProcessedModel {
	var carinfo []apidata.ProcessedModel

	for _, model := range allCarsInfo {
		if model.ManufacturerName == selectedManufacturer && model.CategoryName == selectedCategory {
			carinfo = append(carinfo, model)
		}
	}

	for _, model := range allCarsInfo {
		if model.ManufacturerName != selectedManufacturer && model.CategoryName == selectedCategory || model.ManufacturerName == selectedManufacturer && model.CategoryName != selectedCategory {
			carinfo = append(carinfo, model)
		}
	}
	return carinfo
}

func addCookieFromAManufacturerAndCategoryinfo(w http.ResponseWriter, r *http.Request, selectedManufacturer string, selectedCategory string, allCarsInfo []apidata.ProcessedModel) error {
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
	if err := json.Unmarshal([]byte(string(decodedCookieValue)), &data); err != nil {
		return err
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
