package oauth

import (
	"context"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig = &oauth2.Config{
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

func GetGoogleLoginURL(state string) string {
	return GoogleOauthConfig.AuthCodeURL(state)
}

func ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return GoogleOauthConfig.Exchange(ctx, code)
}
