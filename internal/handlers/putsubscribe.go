package handlers

import (
	"bytes"
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gavink97/gavin-site/internal/components"
	"github.com/gavink97/gavin-site/internal/middleware"
	"github.com/gavink97/gavin-site/internal/store"
)

func PutSubscribeHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*store.User)

	var email string
	if !ok {
		subcookie := getSubscriptionCookie(r)

		emailByte, err := b64.StdEncoding.DecodeString(subcookie.Value)

		if err != nil {
			fmt.Println("error decoding subscription cookie", err)
			return
		}

		email = string(emailByte)
	} else {
		email = user.Email
	}

	firstName := r.FormValue("fname")
	lastName := r.FormValue("lname")
	phoneNumber := r.FormValue("phone")

	member := Member{
		EmailAddress: email,
		Status:       "subscribed",
		MergeFields: map[string]string{
			"FNAME": firstName,
			"LNAME": lastName,
			"PHONE": phoneNumber,
		},
	}

	response, err := updateMemberInfo(member)
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

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<p>Subscription successful for: %s</p>", response.EmailAddress)
}

// https://mailchimp.com/developer/marketing/api/list-members/add-member-to-list/
func updateMemberInfo(member Member) (Member, error) {
	apiServerPrefix := os.Getenv("MAILCHIMP_SERVER_PREFIX")
	apiKey := os.Getenv("MAILCHIMP_API_KEY")
	mailingListID := os.Getenv("MAILCHIMP_LIST_ID")
	hash := createMD5Hash(member.EmailAddress)

	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members/%s?skip_merge_validation=true", apiServerPrefix, mailingListID, hash)

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

func createMD5Hash(string string) string {
	hash := md5.Sum([]byte(strings.ToLower(string)))
	return hex.EncodeToString(hash[:])
}
