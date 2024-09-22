package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func DeleteSubscribeHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	err := deleteMemberFromList(email)
	if err != nil {
		var apiErr APIError
		if jsonErr := json.Unmarshal([]byte(err.Error()), &apiErr); jsonErr == nil {
			fmt.Printf("apiErr Status: %v", apiErr.Status)

			if apiErr.Status == 400 {
				unsubscribingError(w)
			}

			unsubscribingError(w)
		}

		unsubscribingError(w)
	}

	// when testing expect status code 204
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<p>You're now unsubscribed from our newsletter :'(</p>")
}

func deleteMemberFromList(email string) error {
	apiServerPrefix := os.Getenv("MAILCHIMP_SERVER_PREFIX")
	apiKey := os.Getenv("MAILCHIMP_API_KEY")
	mailingListID := os.Getenv("MAILCHIMP_LIST_ID")
	hash := createMD5Hash(email)

	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members/%s", apiServerPrefix, mailingListID, hash)

	req, err := http.NewRequest("DELETE", url, nil)
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
		return errors.New(string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Println(string(body))

	return nil
}

func unsubscribingError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<p>There was an error unsubscribing you from our newsletter. Please contact us to be removed manually or try again later</p>")
}
