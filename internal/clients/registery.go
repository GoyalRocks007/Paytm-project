package clients

import (
	emailclient "paytm-project/internal/clients/email_client"

	"gorm.io/gorm"
)

var (
	registry *ClientRegistry
)

type ClientRegistry struct {
	EmailClient emailclient.IEmailClient
}

func GetRegistry() *ClientRegistry {
	if registry == nil {
		registry = &ClientRegistry{}
	}
	return registry
}

func (cr *ClientRegistry) WithEmailClient(db *gorm.DB) *ClientRegistry {
	cr.EmailClient = emailclient.GetEmailClient(db)
	return cr
}
