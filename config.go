package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	tokenCacheDir := filepath.Join(userConfigDir, "youtube-subscription-feed")
	if err := os.MkdirAll(tokenCacheDir, 0700); err != nil && !errors.Is(err, os.ErrExist) {
		return "", err
	}

	foo := filepath.Join(tokenCacheDir, "token.json")
	return foo, nil
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) error {
	log.Debugf("saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}

	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
