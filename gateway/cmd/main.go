package main

import (
	"auth/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const host = "http://localhost:8801"

func loginUser(httpClient *http.Client, request *models.LoginUserRequest) (*models.LoginUserResponse, error) {
	const uri = "/api/v1/auth/login"

	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(request); err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, strings.Join([]string{host, uri}, ""), body)
	if err != nil {
		return nil, err
	}

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		var errorMsg models.ErrorMessage
		if err := json.NewDecoder(httpResponse.Body).Decode(&errorMsg); err != nil {
			return nil, fmt.Errorf("user not login: status code %d", httpResponse.StatusCode)
		}

		return nil, fmt.Errorf("user not login: status code %d, error: %s", httpResponse.StatusCode, errorMsg.Message)
	}

	response := new(models.LoginUserResponse)
	if err := json.NewDecoder(httpResponse.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("can't unmarshal reponse: %s", err.Error())
	}

	return response, nil
}

func newString(s string) *string {
	return &s
}
func main() {
	httpClient := &http.Client{}

	resp, err := loginUser(httpClient, &models.LoginUserRequest{
		Login:    newString("login"),
		Password: newString("1234"),
	})
	if err != nil {
		log.Printf("login User error: %v", err)
	} else {
		log.Printf("login User: %#v", resp)
	}
}
