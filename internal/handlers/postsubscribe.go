package handlers

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gavink97/gavin-site/internal/components"
)

type APIError struct {
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

type Member struct {
	EmailAddress string            `json:"email_address"`
	Status       string            `json:"status"`
	MergeFields  map[string]string `json:"merge_fields"`
}

// https://medium.com/@reetas/clean-ways-of-adding-new-optional-fields-to-a-golang-struct-99ae2fe9719d
// forgotten email not subscribed
func PostSubscribeHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	status := "subscribed"

	member := Member{
		EmailAddress: email,
		Status:       status,
		MergeFields:  map[string]string{},
	}

	_, err := addMemberToMailingList(member)
	if err != nil {
		var apiErr APIError
		if jsonErr := json.Unmarshal([]byte(err.Error()), &apiErr); jsonErr == nil {
			fmt.Printf("apiErr Status: %v", apiErr.Status)

			if apiErr.Status == 400 {
				w.WriteHeader(http.StatusAccepted)
				c := components.SubscribeExists()
				if err := c.Render(r.Context(), w); err != nil {
					http.Error(w, "Failed to render content", http.StatusInternalServerError)
				}

				fmt.Println(err)
				return
			}

			subscriptionError(w, r)
		}

		subscriptionError(w, r)
	}

	emailCookie := createTemporaryCookie("subscriptionEmail", email)
	http.SetCookie(w, &emailCookie)

	statusCookie := createCookie("subscriptionStatus", status)
	http.SetCookie(w, &statusCookie)

	c := components.SubscribeSuccess()
	if err := c.Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render content", http.StatusInternalServerError)
	}
}

// https://mailchimp.com/developer/marketing/api/list-members/add-member-to-list/
func addMemberToMailingList(member Member) (Member, error) {
	apiServerPrefix := os.Getenv("MAILCHIMP_SERVER_PREFIX")
	apiKey := os.Getenv("MAILCHIMP_API_KEY")
	mailingListID := os.Getenv("MAILCHIMP_LIST_ID")

	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members?skip_merge_validation=true", apiServerPrefix, mailingListID)

	jsonData, err := json.Marshal(member)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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

func subscriptionError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	c := components.SubscribeError()
	if err := c.Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render content", http.StatusInternalServerError)
	}
}

func createTemporaryCookie(key, value string) http.Cookie {
	cookieValue := b64.StdEncoding.EncodeToString([]byte(value))

	expiration := time.Now().Add(15 * time.Minute)
	cookie := http.Cookie{
		Name:     key,
		Value:    cookieValue,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	return cookie
}

func createCookie(key, value string) http.Cookie {
	cookieValue := b64.StdEncoding.EncodeToString([]byte(value))

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     key,
		Value:    cookieValue,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	return cookie
}

func getSubscriptionCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("subscriptionEmail")

	if err != nil {
		fmt.Println("No cookie is present", err)
		return nil
	}
	return cookie
}

func checkStatusCookie(r *http.Request) bool {
	_, err := r.Cookie("subscriptionStatus")

	if err != nil {
		fmt.Println("No cookie is present", err)
		return false
	}
	return true
}
