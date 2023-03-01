package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Yothgewalt/aufruf-jaeger-bot/module/entity"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	_TOKEN_FILENAME string = "token.json"
)

var (
	logger zerolog.Logger
)

type googleClient struct {
	logger zerolog.Logger
}

// Initial a function for retrieve two (such as logger, file that contains token for authorize etc.,)
func NewGoogleClientRepository(loggerSource zerolog.Logger) entity.GoogleClient {
	logger = loggerSource

	return &googleClient{
		logger: loggerSource,
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func (gc *googleClient) GetClient(jsonKey []byte, scope ...string) *http.Client {
	config, err := google.ConfigFromJSON(jsonKey, scope...)
	if err != nil {
		gc.logger.Err(err).Msg("Unable to parse client secret file to config")
	}

	tokFile := _TOKEN_FILENAME
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		logger.Fatal().Err(err).Msg("Unable to read authorization code")
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to retrieve token from web")
	}

	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	logger.Info().Str("Path", path).Msg("Saving credential file to")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to cache oauth token")
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
