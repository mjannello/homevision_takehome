package web

import (
	"encoding/json"
	"fmt"
	http2 "homevision/pkg/http"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const homevisionURL = "http://app-homevision-staging.herokuapp.com/api_project/houses"

// TODO: create a model.go to include model structs

type House struct {
	ID        int    `json:"id"`
	Address   string `json:"address"`
	Homeowner string `json:"homeowner"`
	Price     int    `json:"price"`
	PhotoURL  string `json:"photoURL"`
}

type HousesResponse struct {
	Houses []House `json:"houses"`
	OK     bool    `json:"ok"`
}

type Image struct {
	ID      int
	Address string
	Content string
	Name    string
}

func FetchHousesInfo(totalPages, perPage int, httpClient http2.RetryableHTTPClient) ([]House, error) {
	var housesList []House
	var wg sync.WaitGroup
	responsesChan := make(chan *http.Response, totalPages)
	// TODO: Add an errChan to catch possible errors
	for i := 1; i <= totalPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()

			url := fmt.Sprintf("%s?page=%d&per_page=%d", homevisionURL, page, perPage)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}

			resp, err := httpClient.Do(req)
			if err != nil {
				fmt.Println("Error on request:", err)
				return
			}

			responsesChan <- resp
		}(i)
	}

	go func() {
		wg.Wait()
		close(responsesChan)
	}()

	for resp := range responsesChan {
		var housesResponse HousesResponse
		err := json.NewDecoder(resp.Body).Decode(&housesResponse)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return nil, err
		}

		if housesResponse.OK {
			housesList = append(housesList, housesResponse.Houses...)
		}
		resp.Body.Close()
	}

	return housesList, nil
}

func DownloadImageContent(url string, client http2.RetryableHTTPClient) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for image content: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download image content: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image content, status code: %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image content: %v", err)
	}

	return string(content), nil
}

func ProcessHouseImages(houses []House, client http2.RetryableHTTPClient) {
	var wg sync.WaitGroup
	imagesChan := make(chan Image, len(houses))

	wg.Add(1)
	go func() {
		defer wg.Done()

		for image := range imagesChan {
			fmt.Printf("Processing image %s\n", image.Name)
			content, err := DownloadImageContent(image.Content, client)
			if err != nil {
				fmt.Printf("Error downloading image content: %v\n", err)
				continue
			}

			ext := filepath.Ext(image.Content)
			fileName := fmt.Sprintf("%d-%s%s", image.ID, cleanFileName(image.Address), ext)
			// TODO: add a custom filepath. Make it a param for the current method
			filePath := filepath.Join(fileName)

			err = os.WriteFile(filePath, []byte(content), 0644)
			if err != nil {
				fmt.Printf("Error saving image %s: %v\n", fileName, err)
			} else {
				fmt.Printf("Image saved successfully: %s\n", fileName)
			}
		}
	}()

	for _, house := range houses {
		image := Image{ID: house.ID, Address: house.Address, Content: house.PhotoURL, Name: fmt.Sprintf("%d-%s", house.ID, cleanFileName(house.Address))}
		imagesChan <- image
	}
	close(imagesChan)

	wg.Wait()
}

func cleanFileName(s string) string {
	// TODO: use regex
	result := strings.ReplaceAll(s, " ", "_")
	result = strings.ReplaceAll(result, ",", "")
	result = strings.ReplaceAll(result, ".", "")
	return result
}
