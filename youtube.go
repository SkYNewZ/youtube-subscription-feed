package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Debugf("Go to the following link in your browser then type the authorization code:\n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}

	return tok, nil
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var t oauth2.Token
	return &t, json.NewDecoder(f).Decode(&t)
}

// GetYouTubeClient uses a context.Context and oauth2.Config to retrieve a oauth2.Token
// then generate a Client. It returns the generated Client.
func GetYouTubeClient(ctx context.Context) (*youtube.Service, error) {
	b, err := ioutil.ReadFile(secretFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	// obtain the token config
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	// read the token from filesystem
	cacheFile, err := tokenCacheFile()
	if err != nil {
		return nil, fmt.Errorf("unable to get path to cached credential file: %v", err)
	}

	log.Debugf("reading saved token from: %s", cacheFile)
	token, err := tokenFromFile(cacheFile)

	if err != nil {
		log.Debugf("error reading token from file: %v", err)

		token, err = getTokenFromWeb(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("unable to authenticate: %v", err)
		}

		if err := saveToken(cacheFile, token); err != nil {
			log.Warningf("cannot save token to filesystem: %v", err)
		}
	}

	return youtube.NewService(
		ctx,
		option.WithScopes(youtube.YoutubeReadonlyScope),
		option.WithTokenSource(config.TokenSource(ctx, token)),
	)
}

// listSubscriptionsChannelIDs return a channel with each current user subscription
func listSubscriptionsChannelIDs(ctx context.Context, service *youtube.Service) ([]*youtube.Subscription, error) {
	var results = make([]*youtube.Subscription, 0)
	err := service.Subscriptions.
		List([]string{youtubePartSnippet}).
		Mine(true).
		MaxResults(20).
		Order("alphabetical").
		Pages(ctx, func(response *youtube.SubscriptionListResponse) error {
			results = append(results, response.Items...)
			return nil
		})

	return results, err
}
