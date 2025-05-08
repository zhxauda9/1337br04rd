package externalapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"1337b04rd/internal/domain"
)

type RickAndMortyClient struct {
	baseURL string
}

type apiResponse struct {
	Info struct {
		Next string `json:"next"`
	} `json:"info"`
	Results []domain.Character `json:"results"`
}

func NewRickAndMortyClient() *RickAndMortyClient {
	return &RickAndMortyClient{
		baseURL: "https://rickandmortyapi.com/api",
	}
}

// FetchAllCharacters fetches all characters from the Rick and Morty API
func (c *RickAndMortyClient) FetchAllCharacters() ([]domain.Character, error) {
	var allCharacters []domain.Character
	nextURL := fmt.Sprintf("%s/character", c.baseURL)

	for nextURL != "" {
		resp, err := http.Get(nextURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch characters: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var apiResp apiResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		allCharacters = append(allCharacters, apiResp.Results...)
		nextURL = apiResp.Info.Next
	}

	return allCharacters, nil
}

func (c *RickAndMortyClient) FetchRandomCharacter() (*domain.Character, error) {
	randomID := rand.Intn(826) + 1
	resp, err := http.Get(fmt.Sprintf("%s%d", c.baseURL, randomID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch character: %w", err)
	}
	defer resp.Body.Close()

	var char domain.Character
	if err := json.NewDecoder(resp.Body).Decode(&char); err != nil {
		return nil, fmt.Errorf("failed to decode character: %w", err)
	}

	return &char, nil
}
