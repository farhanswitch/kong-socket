package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AuthProvider struct {
	authURL string
}

var provider *AuthProvider

type ResponseDetail struct {
	Id     uint16 `json:"id"`
	Role   uint16 `json:"role"`
	Status uint16 `json:"status"`
	Iat    uint32 `json:"iat"`
	Exp    uint32 `json:"exp"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	House  string `json:"house"`
}
type DataResponse struct {
	Message string         `json:"statusCode"`
	Data    ResponseDetail `json:"data"`
	Status  uint16         `json:"status"`
}

func (ap AuthProvider) AuthRequest(token string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, ap.authURL, nil)
	if err != nil {

		log.Printf("Error: %v\n", err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	// Lakukan HTTP Request ke service Autentikasi
	res, err := client.Do(req)
	if err != nil {

		log.Printf("Error: %v\n", err)
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var accessResponse DataResponse
	err = json.Unmarshal(body, &accessResponse)
	if err != nil {
		return "", err
	}
	if accessResponse.Status != http.StatusOK {
		return "", fmt.Errorf("%s", "Access Forbidden")
	}
	var room string = "sendPriceU"
	if accessResponse.Data.Type == "cln" {
		if accessResponse.Data.House == "ALLR9999" {
			room = "sendPriceR"
		}
	}
	return room, nil
}
func AuthProviderFactory(authURL string) AuthProvider {
	if provider == nil {
		provider = &AuthProvider{authURL}
	}
	return *provider
}
