package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gavink97/gavin-site/internal/middleware"
	"github.com/gavink97/gavin-site/internal/store"
)

func GetSubscribeHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*store.User)
	if !ok {
		fmt.Println("user is not authenticated")
	}

	query := r.URL.Query()

	var email string
	if len(query) > 0 {
		email = query.Get("email")
	} else {
		email = user.Email
	}

	hash := createMD5Hash(email)
	response, err := GetMemberInfo(hash)
	if err != nil {
		var apiErr APIError
		if jsonErr := json.Unmarshal([]byte(err.Error()), &apiErr); jsonErr == nil {
			fmt.Printf("apiErr Status: %v", apiErr.Status)

			getMemberInfoError(w)
		}

		getMemberInfoError(w)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(string(response))
	if err != nil {
		return
	}
}

func GetMemberInfo(hash string) ([]byte, error) {
	apiServerPrefix := os.Getenv("MAILCHIMP_SERVER_PREFIX")
	apiKey := os.Getenv("MAILCHIMP_API_KEY")
	mailingListID := os.Getenv("MAILCHIMP_LIST_ID")

	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members/%s", apiServerPrefix, mailingListID, hash)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Body)
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Println(string(body))

	return body, nil
}

func getMemberInfoError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := "An error occured when fetching member data"
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		return
	}
}
