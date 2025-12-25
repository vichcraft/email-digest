package email

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context, credPath string, tokenPath string) (*gmail.Service, error) {

	//Read credentials.json
	credBytes, err := os.ReadFile(credPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to read credentials file: %v", err)
	}

	//Configure OAuth2
	config, err := google.ConfigFromJSON(credBytes, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse credentials.json: %v", err)
	}

	//Get OAuth2 token
	token, err := getToken(config, tokenPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to get token: %v", err)
	}

	//Initialize authenticated service
	client := config.Client(ctx, token)
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("Unable to create Gmail service: %v", err)
	}

	return service, nil
}

func getToken(config *oauth2.Config, tokenPath string) (*oauth2.Token, error) {
	token, err := loadToken(tokenPath)
	if err == nil {
		return token, nil
	}

	token, err = getTokenFromWeb(config)
	if err != nil {
		return nil, err
	}

	if err := saveToken(tokenPath, token); err != nil {
		log.Printf("Unable to save token: %v", err)
	}

	return token, nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authUrl := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to this URL in your browser:\n%v\n\n", authUrl)
	fmt.Print("Enter the authorization code: ")

	var authCode string

	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve token: %v", err)
	}

	return token, nil

}

func loadToken(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

func saveToken(path string, token *oauth2.Token) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}
