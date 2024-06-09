package usecases

import (
	"context"
	config "github.com/Nixonxp/discord/auth/configs"
	"github.com/Nixonxp/discord/auth/internal/app/models"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	pkgerrors "github.com/Nixonxp/discord/auth/pkg/errors"
	log "github.com/Nixonxp/discord/auth/pkg/logger"
	"github.com/Nixonxp/discord/auth/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Deps struct {
	UserService usecases.UserServiceInterface
	Log         *log.Logger
	OauthSvc    usecases.OAuthServiceInterface
	Cfg         *config.Config
}

type AuthUsecase struct {
	Deps
}

var _ usecases.UsecaseInterface = (*AuthUsecase)(nil)

func NewAuthUsecase(d Deps) usecases.UsecaseInterface {
	return &AuthUsecase{
		Deps: d,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, registerInfo usecases.RegisterUserInfo) (*models.User, error) {
	password := []byte(registerInfo.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return &models.User{}, pkgerrors.Wrap("password hashing error", err)
	}

	registerInfo.Password = string(hashedPassword)

	user, err := u.UserService.Register(ctx, registerInfo)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, loginInfo usecases.LoginUserInfo) (*models.LoginResult, error) {
	user, err := u.UserService.GetUserForLogin(ctx, loginInfo)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
	if err != nil {
		return nil, pkgerrors.Wrap("wrong password or login", models.ErrCredInvalid)
	}

	token, refreshToken, err := u.generateTokensForUserId(user.UserID.String())
	if err != nil {
		return nil, err
	}

	return &models.LoginResult{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (u *AuthUsecase) generateTokensForUserId(userId string) (string, string, error) {
	token, err := generateJWTToken(u.Cfg.Application.AuthSecretKey, userId)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateRefreshJWTToken(u.Cfg.Application.AuthRefreshSecretKey, userId)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func (u *AuthUsecase) Refresh(_ context.Context, refreshToken string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.Cfg.Application.AuthRefreshSecretKey), nil
	})
	if err != nil {
		return "", pkgerrors.Wrap("Invalid or expired JWT token", models.Unauthenticated)
	}

	if !token.Valid {
		return "", pkgerrors.Wrap("Invalid JWT token", models.Unauthenticated)
	}

	claims := token.Claims.(*jwt.StandardClaims)
	userID := claims.Id

	newToken, err := generateJWTToken(u.Cfg.Application.AuthSecretKey, userID)
	if !token.Valid {
		return "", pkgerrors.Wrap("error generate JWT token", models.Unauthenticated)
	}

	return newToken, nil
}

func generateJWTToken(secretKey string, userID string) (string, error) {
	// Define the expiration time for the token.
	expirationTime := time.Now().Add(1 * time.Hour)
	// Create the JWT claims, which include the user ID and expiration time.
	claims := jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Id:        userID,
	}
	// Create the JWT token with the claims and a secret key.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func generateRefreshJWTToken(refreshSecretKey string, userID string) (string, error) {
	// Define the expiration time for the token.
	expirationTime := time.Now().Add(24 * time.Hour)
	// Create the JWT claims, which include the user ID and expiration time.
	claims := jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Id:        userID,
	}
	// Create the JWT token with the claims and a secret key.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(refreshSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (u *AuthUsecase) OauthLogin(_ context.Context, _ usecases.OauthLoginRequest) (*models.OauthLoginResult, error) {
	url := u.OauthSvc.AuthCodeURL("state")
	return &models.OauthLoginResult{
		Code: url,
	}, nil
}

func (u *AuthUsecase) OauthLoginCallback(ctx context.Context, req usecases.OauthLoginCallbackRequest) (*models.LoginResult, error) {
	if req.State != "state" {
		return nil, pkgerrors.Wrap("Invalid state parameter", models.Unauthenticated)
	}

	userInfo, err := u.OauthSvc.ExchangeClient(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(utils.GenerateRandomString(10)), bcrypt.MinCost)
	if err != nil {
		return nil, pkgerrors.Wrap("password hashing error", err)
	}

	user, err := u.UserService.CreateOrCreateUser(ctx, usecases.GetOrCreateUserRequest{
		Login:          userInfo.Email,
		Name:           userInfo.Name,
		Email:          userInfo.Email,
		Password:       string(hashedPassword),
		AvatarPhotoUrl: userInfo.AvatarPhotoUrl,
		OauthId:        userInfo.OauthId,
	})
	if err != nil {
		return nil, err
	}

	authToken, refreshToken, err := u.generateTokensForUserId(user.UserID.String())
	if err != nil {
		return nil, err
	}

	return &models.LoginResult{
		Token:        authToken,
		RefreshToken: refreshToken,
	}, nil
}
