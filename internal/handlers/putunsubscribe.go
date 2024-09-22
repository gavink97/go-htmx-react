package handlers

import (
	"bytes"
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

func PutUnsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*store.User)
	if !ok {
		fmt.Println("user is not authenticated")
	}

	query := r.URL.Query()

	formValue := r.FormValue("email")

	var email string
	if len(query) > 0 {
		// send something that can be decoded by server
		email = query.Get("hash")

		// decode the value
	} else if user != nil {
		email = user.Email
	} else if formValue != "" {
		email = formValue
	} else {
		fmt.Println("there was an error getting the email address")
		return
	}

	member := Member{
		EmailAddress: email,
		Status:       "unsubscribed",
		MergeFields:  map[string]string{},
	}

	_, err := unsubscribeMember(member)
	if err != nil {
		var apiErr APIError
		if jsonErr := json.Unmarshal([]byte(err.Error()), &apiErr); jsonErr == nil {
			fmt.Printf("apiErr Status: %v", apiErr.Status)

			unsubscribingError(w)
		}

		unsubscribingError(w)
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<p>You've successfully unsubscribed to our newsletter :'(</p>")
}

func unsubscribeMember(member Member) (Member, error) {
	apiServerPrefix := os.Getenv("MAILCHIMP_SERVER_PREFIX")
	apiKey := os.Getenv("MAILCHIMP_API_KEY")
	mailingListID := os.Getenv("MAILCHIMP_LIST_ID")

	hash := createMD5Hash(member.EmailAddress)

	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members/%s", apiServerPrefix, mailingListID, hash)

	jsonData, err := json.Marshal(member)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
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
		return Member{}, errors.New(string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Println(string(body))

	return member, nil
}
