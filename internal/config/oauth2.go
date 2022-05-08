package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func SetupConfigGoogle() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     "6116145082-n6bu7lpemg1cicrooa19gepmmhh9n4uu.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8081/hello",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}
