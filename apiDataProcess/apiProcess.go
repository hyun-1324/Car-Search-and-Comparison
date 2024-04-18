/*
This file handles the retrieval and processing of data from an API server.
It fetches car models, manufacturers, and categories, and then transforms
this data into a structured format suitable for use within this project.
*/

package apidata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

const ApiURL = "http://localhost:3000/"

// ModelApi holds car model information retrieved from the API server.
type ModelApi struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	ManufacturerId int    `json:"manufacturerId"`
	CategoryId     int    `json:"categoryId"`
	Year           int    `json:"year"`
	Specifications struct {
		Engine       string `json:"engine"`
		Horsepower   int    `json:"horsepower"`
		Transmission string `json:"transmission"`
		Drivetrain   string `json:"drivetrain"`
	} `json:"specifications"`
	Image string `json:"image"`
}

// ManufacturerApi holds manufacturer information retrieved from the API server.
type ManufacturerApi struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

// CategoryApi holds category information retrieved from the API server.
type CategoryApi struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// ProcessedModel holds the structured data derived from the API.
type ProcessedModel struct {
	Id                       int    `json:"id"`
	Name                     string `json:"name"`
	ManufacturerName         string `json:"manufacturerName"`
	ManufacturerCountry      string `json:"manufacturerCountry"`
	ManufacturerFoundingYear int    `json:"manufacturerFoundingYear"`
	CategoryName             string `json:"categoryName"`
	Year                     int    `json:"year"`
	Specifications           struct {
		Engine       string `json:"engine"`
		Horsepower   int    `json:"horsepower"`
		Transmission string `json:"transmission"`
		Drivetrain   string `json:"drivetrain"`
	} `json:"specifications"`
	Image string `json:"image"`
}

// Converts API data to a processed format used in this project.
func ProcessedApiData() ([]ProcessedModel, error) {
	var wg sync.WaitGroup
	results := make(chan []ProcessedModel)
	errors := make(chan error)

	var models []ModelApi
	var manufacturers []ManufacturerApi
	var categories []CategoryApi

	wg.Add(3)

	go func() {
		defer wg.Done()
		var err error
		models, err = fetchModels()
		if err != nil {
			errors <- fmt.Errorf("fetchModel error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		manufacturers, err = fetchManufacturers()
		if err != nil {
			errors <- fmt.Errorf("fetchManufacturer error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		categories, err = fetchCategories()
		if err != nil {
			errors <- fmt.Errorf("fetchCategory error: %v", err)
		}
	}()

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	for err := range errors {
		return nil, err
	}

	var processedModels []ProcessedModel
	for _, model := range models {
		newModel := ProcessedModel{
			Id:   model.Id,
			Name: model.Name,
			Year: model.Year,
			Specifications: struct {
				Engine       string `json:"engine"`
				Horsepower   int    `json:"horsepower"`
				Transmission string `json:"transmission"`
				Drivetrain   string `json:"drivetrain"`
			}{
				Engine:       model.Specifications.Engine,
				Horsepower:   model.Specifications.Horsepower,
				Transmission: model.Specifications.Transmission,
				Drivetrain:   model.Specifications.Drivetrain,
			},
			Image: model.Image,
		}

		for _, manufacturer := range manufacturers {
			if model.ManufacturerId == manufacturer.Id {
				newModel.ManufacturerName = manufacturer.Name
				newModel.ManufacturerCountry = manufacturer.Country
				newModel.ManufacturerFoundingYear = manufacturer.FoundingYear
				break
			}
		}
		for _, category := range categories {
			if model.CategoryId == category.Id {
				newModel.CategoryName = category.Name
				break
			}
		}
		processedModels = append(processedModels, newModel)

	}
	return processedModels, nil
}

// fetches model information from the API server
func fetchModels() ([]ModelApi, error) {
	url := ApiURL + "api/models"

	responseData, err := getDataFromApi(url)
	if err != nil {
		return nil, err
	}

	var modelList []ModelApi
	if err := json.Unmarshal(responseData, &modelList); err != nil {
		return nil, err
	}

	return modelList, nil

}

// fetches Manufacturer information from the API server
func fetchManufacturers() ([]ManufacturerApi, error) {
	url := ApiURL + "api/manufacturers"

	responseData, err := getDataFromApi(url)
	if err != nil {
		return nil, err
	}

	var manufacturerList []ManufacturerApi
	if err := json.Unmarshal(responseData, &manufacturerList); err != nil {
		return nil, err
	}

	return manufacturerList, nil

}

// fetches categories information from the API server
func fetchCategories() ([]CategoryApi, error) {
	url := ApiURL + "api/categories"

	responseData, err := getDataFromApi(url)
	if err != nil {
		return nil, err
	}

	var categoryList []CategoryApi
	if err := json.Unmarshal(responseData, &categoryList); err != nil {
		fmt.Println("fetchcategory unmarshal error")
		return nil, err
	}

	return categoryList, nil

}

// Retrieves unprocessed data from the API server.
func getDataFromApi(url string) (responseData []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-200 status: %d", resp.StatusCode)
	}

	responseData, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil

}
