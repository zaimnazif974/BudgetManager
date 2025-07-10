package config

import (
	// "github.com/gorilla/sessions"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"

	// "github.com/markbates/goth/gothic"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type GoogleAuth struct {
	ClientID     string
	ClientSecret string
}

type SessionConfig struct {
	Key    string
	MaxAge int32
	IsProd bool
}

func LoadGoogleAuth() *GoogleAuth {
	return &GoogleAuth{
		ClientID:     GetEnv("CLIENT_ID", "nil"),
		ClientSecret: GetEnv("CLIENT_SECRET", "nil"),
	}
}

func LoadSessionConfig() *SessionConfig {
	return &SessionConfig{
		Key:    GetEnv("SESSION_SECRET", "nil"),
		MaxAge: 86400 * 30,
		IsProd: false,
	}
}

func AuthConfig() {

	//Config
	googleConfig := LoadGoogleAuth()
	sessionConfig := LoadSessionConfig()

	//Session Store
	store := sessions.NewCookieStore([]byte(sessionConfig.Key))
	store.MaxAge(int(sessionConfig.MaxAge))

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = sessionConfig.IsProd
	store.Options.SameSite = http.SameSiteLaxMode

	gothic.Store = store

	//oAuth
	goth.UseProviders(
		google.New(
			googleConfig.ClientID,
			googleConfig.ClientSecret,
			"http://localhost:8080/api/v1/auth/google/callback",
			"email", "profile",
		),
	)
}
