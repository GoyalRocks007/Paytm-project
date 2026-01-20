package gmailclient

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailClient struct {
	GmailRepo IGmailClientRepo
}

func (gc *GmailClient) StartAuth(ctx context.Context) (string, error) {
	config, err := GetOAuthConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get OAuth config: %v", err)
	}

	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return url, nil
}

func (gc *GmailClient) HandleCallback(ctx context.Context, code string) error {
	config, err := GetOAuthConfig()
	if err != nil {
		return err
	}

	// Exchange code for token
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return err
	}

	// Get user ID from query params - for demo, using a fixed user
	userID := os.Getenv("GOOGLE_EMAIL_ADDRESS")

	// Save token to database
	err = gc.GmailRepo.SaveTokenToDB(userID, token)
	if err != nil {
		return err
	}

	return nil
}

func (gc *GmailClient) FetchToken(ctx context.Context, config *oauth2.Config, userID string) (*oauth2.Token, error) {
	// Load saved token from database
	token, err := gc.GmailRepo.LoadTokenFromDB(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load token: %v", err)
	}

	// Check if token needs refresh
	tokenSource := config.TokenSource(ctx, token)
	freshToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %v", err)
	}

	// If token was refreshed, save the new one to database
	if freshToken.AccessToken != token.AccessToken {
		err = gc.GmailRepo.SaveTokenToDB(userID, freshToken)
		if err != nil {
			log.Printf("Warning: failed to save refreshed token: %v", err)
		}
	}
	return freshToken, nil
}

func (gc *GmailClient) SendHtmlEmail(ctx context.Context, userID, to, subject, body string) error {
	config, err := GetOAuthConfig()
	if err != nil {
		return fmt.Errorf("failed to get OAuth config: %v", err)
	}

	freshToken, err := gc.FetchToken(ctx, config, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch token: %v", err)
	}

	client := config.Client(ctx, freshToken)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create Gmail service: %v", err)
	}

	msgStr := fmt.Sprintf(
		"From: me\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
			"%s",
		to, subject, body,
	)

	message := &gmail.Message{
		Raw: EncodeMessage(msgStr),
	}

	_, err = srv.Users.Messages.Send("me", message).Do()
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (gc *GmailClient) SendEmail(ctx context.Context, userID, to, subject, body string) error {

	// Get OAuth config
	config, err := GetOAuthConfig()
	if err != nil {
		return fmt.Errorf("failed to get OAuth config: %v", err)
	}

	freshToken, err := gc.FetchToken(ctx, config, userID)

	if err != nil {
		return fmt.Errorf("failed to fetch token: %v", err)
	}

	// Create authenticated client
	client := config.Client(ctx, freshToken)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create Gmail service: %v", err)
	}

	// Prepare and send email
	msgStr := fmt.Sprintf("From: me\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)
	message := &gmail.Message{
		Raw: EncodeMessage(msgStr),
	}

	_, err = srv.Users.Messages.Send("me", message).Do()
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
