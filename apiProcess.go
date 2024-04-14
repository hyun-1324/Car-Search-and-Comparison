package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

type ManufacturerApi struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

type CategoryApi struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ProcessedModel struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	ManufacturerName string `json:"manufacturerName"`
	CategoryName     string `json:"categoryName"`
	Year             int    `json:"year"`
	Specifications   struct {
		Engine       string `json:"engine"`
		Horsepower   int    `json:"horsepower"`
		Transmission string `json:"transmission"`
		Drivetrain   string `json:"drivetrain"`
	} `json:"specifications"`
	Image string `json:"image"`
}

func processedApiData() ([]ProcessedModel, error) {
	var processedModels []ProcessedModel
	models, err := fetchModels()
	if err != nil {
		return nil, err
	}
	manufatures, err := fetchManufacturers()
	if err != nil {
		return nil, err
	}
	categories, err := fetchCategory()
	if err != nil {
		return nil, err
	}
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

		for _, manufacturer := range manufatures {
			if model.ManufacturerId == manufacturer.Id {
				newModel.ManufacturerName = manufacturer.Name
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

func fetchModels() ([]ModelApi, error) {
	url := "http://localhost:3000/api/models"
	client := &http.Client{}

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

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var modelList []ModelApi
	if err := json.Unmarshal(responseData, &modelList); err != nil {
		return nil, err
	}

	return modelList, nil

}

func fetchManufacturers() ([]ManufacturerApi, error) {
	url := "http://localhost:3000/api/manufacturers"
	client := &http.Client{}

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

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var manufacturerList []ManufacturerApi
	if err := json.Unmarshal(responseData, &manufacturerList); err != nil {
		return nil, err
	}

	return manufacturerList, nil

}

func fetchCategory() ([]CategoryApi, error) {
	url := "https://localhost:3000/api/categories"
	client := &http.Client{}

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

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var categoryList []CategoryApi
	if err := json.Unmarshal(responseData, &categoryList); err != nil {
		return nil, err
	}

	return categoryList, nil

}
