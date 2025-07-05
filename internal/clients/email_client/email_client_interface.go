package emailclient

import (
	"context"

	gmailclient "paytm-project/internal/clients/email_client/gmail_client"
	"paytm-project/internal/models"

	"gorm.io/gorm"
)

var (
	gmailClientInstance IEmailClient
	GetEmailClient      = func(db *gorm.DB) IEmailClient {
		if gmailClientInstance != nil {
			return gmailClientInstance
		}
		return initEmailClient(db)
	}
	initEmailClient = func(db *gorm.DB) IEmailClient {
		gmailClientInstance = &gmailclient.GmailClient{
			GmailRepo: &gmailclient.GmailClientRepo{
				BaseRepo: models.BaseRepo{
					Db: db,
				},
			},
		}
		return gmailClientInstance
	}
)

type IEmailClient interface {
	SendEmail(ctx context.Context, userID, to, subject, body string) error
	StartAuth(ctx context.Context) (string, error)
	HandleCallback(ctx context.Context, code string) error
}
