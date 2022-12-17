package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	client = &http.Client{
		Timeout: 2 * time.Second,
	}

	BaseURL = "https://shazam.p.rapidapi.com"
)

func makeRequest(method, url string, body interface{}) error {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", BaseURL, url), nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-RapidAPI-Key", os.Getenv("RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", os.Getenv("RAPID_API_HOST"))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, body)
	if err != nil {
		return err
	}

	return nil
}

func get(url string, body interface{}) {
	makeRequest(http.MethodGet, url, body)
}
