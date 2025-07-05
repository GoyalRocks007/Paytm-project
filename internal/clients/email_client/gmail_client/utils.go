package gmailclient

import (
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func GetOAuthConfig() (*oauth2.Config, error) {
	credentialsJSON := os.Getenv("GOOGLE_CREDENTIALS_JSON")
	if credentialsJSON == "" {
		return nil, fmt.Errorf("GOOGLE_CREDENTIALS_JSON environment variable not set")
	}

	return google.ConfigFromJSON([]byte(credentialsJSON), gmail.GmailSendScope)
}

func EncodeMessage(message string) string {
	return base64.URLEncoding.EncodeToString([]byte(message))
}
