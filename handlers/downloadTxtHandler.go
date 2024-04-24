package handlers

import (
	apidata "cars/apiDataProcess"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func DownloadTxt(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form request: ", err)
		http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
		return
	}
	carName := r.FormValue("downloadtxt")
	acarInfo, err := findACarInfoByOnlyName(carName)
	if err != nil {
		log.Println("Error fetching data: ", err)
		http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
		return
	}

	data := fmt.Sprintf(`%v
Basic information:
	Category: %v
	Production year: %v
Manufacturer's information:
	Manufacturer: %v
	Home country: %v
	Founding year: %v
Specifications:
	Engine: %v
	Horsepower: %v
	Transmission: %v
	Drivetrain: %v
	`, acarInfo.Name, acarInfo.CategoryName, acarInfo.Year,
		acarInfo.ManufacturerName, acarInfo.ManufacturerCountry, acarInfo.ManufacturerFoundingYear,
		acarInfo.Specifications.Engine, acarInfo.Specifications.Horsepower,
		acarInfo.Specifications.Transmission, acarInfo.Specifications.Drivetrain)

	carnameSpaceReplaced := strings.ReplaceAll(acarInfo.Name, " ", "_")
	carnameLowLetters := strings.ToLower(carnameSpaceReplaced)

	filename := fmt.Sprintf("%v.txt", carnameLowLetters)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))
	w.Header().Set("Content-Type", "text/plain")

	err = os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		log.Println("Error creating a text file", err)
		http.Error(w, "Sorry, something went wrong on our end. We're working to fix it!", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, filename)

	defer func() {
		err := os.Remove(filename)
		if err != nil {
			log.Println("Error deleting a text file", err)
		}
	}()

}

func findACarInfoByOnlyName(carName string) (apidata.ProcessedModel, error) {
	var carInfo apidata.ProcessedModel
	allCarsInfo, err := apidata.ProcessedApiData()
	if err != nil {
		return apidata.ProcessedModel{}, err
	}

	for _, model := range allCarsInfo {
		if model.Name == carName {
			carInfo = model
		}
	}

	return carInfo, nil
}
