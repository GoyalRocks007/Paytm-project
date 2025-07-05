package gmailclient

import (
	"fmt"
	"paytm-project/internal/models"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type IGmailClientRepo interface {
	SaveTokenToDB(userID string, token *oauth2.Token) error
	LoadTokenFromDB(userID string) (*oauth2.Token, error)
}

type GmailClientRepo struct {
	models.BaseRepo
}

func (gr *GmailClientRepo) SaveTokenToDB(userID string, token *oauth2.Token) error {
	tokenStorage := TokenStorage{
		UserID:       userID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
	}

	// Use UPSERT - update if exists, create if not
	result := gr.Db.Where("user_id = ?", userID).FirstOrCreate(&tokenStorage)
	if result.Error != nil {
		return fmt.Errorf("failed to save token: %v", result.Error)
	}

	// If record existed, update it
	if result.RowsAffected == 0 {
		tokenStorage.AccessToken = token.AccessToken
		tokenStorage.RefreshToken = token.RefreshToken
		tokenStorage.TokenType = token.TokenType
		tokenStorage.Expiry = token.Expiry

		result = gr.Db.Save(&tokenStorage)
		if result.Error != nil {
			return fmt.Errorf("failed to update token: %v", result.Error)
		}
	}

	return nil
}

func (gr *GmailClientRepo) LoadTokenFromDB(userID string) (*oauth2.Token, error) {
	var tokenStorage TokenStorage

	result := gr.Db.Where("user_id = ?", userID).First(&tokenStorage)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no token found for user %s", userID)
		}
		return nil, fmt.Errorf("failed to load token: %v", result.Error)
	}

	token := &oauth2.Token{
		AccessToken:  tokenStorage.AccessToken,
		RefreshToken: tokenStorage.RefreshToken,
		TokenType:    tokenStorage.TokenType,
		Expiry:       tokenStorage.Expiry,
	}

	return token, nil
}
