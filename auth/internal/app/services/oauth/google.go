package services

import (
	"context"
	"encoding/json"
	config "github.com/Nixonxp/discord/auth/configs"
	"github.com/Nixonxp/discord/auth/internal/app/models"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	pkgerrors "github.com/Nixonxp/discord/auth/pkg/errors"
	logger "github.com/Nixonxp/discord/auth/pkg/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
)

type GoogleOauth struct {
	config oauth2.Config
	log    *logger.Logger
}

func NewGoogleOAuth(cfg *config.Config, log *logger.Logger) *GoogleOauth {
	oauthConfig := oauth2.Config{
		ClientID:     cfg.Application.OAuthClientID,
		ClientSecret: cfg.Application.OAuthClientSecret,
		RedirectURL:  cfg.Application.OAuthRedirectUrl,
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	return &GoogleOauth{
		config: oauthConfig,
		log:    log,
	}
}

func (g *GoogleOauth) AuthCodeURL(state string) string {
	return g.config.AuthCodeURL(state)
}

func (g *GoogleOauth) ExchangeClient(ctx context.Context, code string) (*usecases.UserInfo, error) {
	token, err := g.config.Exchange(context.Background(), code)
	if err != nil {
		g.log.WithContext(ctx).WithError(err).WithField("token", token).Error("Failed to exchange token")
		return nil, pkgerrors.Wrap("Failed to exchange token", models.Unauthenticated)
	}

	client := g.config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		g.log.WithContext(ctx).WithError(err).WithField("token", token).Error("fail auth client from token")
		return nil, pkgerrors.Wrap("Failed to get user info", models.Unauthenticated)
	}

	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	var userInfo usecases.UserInfo
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		g.log.WithContext(ctx).WithError(err).WithField("token", token).Error("Failed unpack user info")
		return nil, pkgerrors.Wrap("Failed unpack user info", err)
	}

	return &userInfo, nil
}
