package usecases

import (
	"context"
	"encoding/json"
	"github.com/Nixonxp/discord/auth/internal/app/models"
	pkgerrors "github.com/Nixonxp/discord/auth/pkg/errors"
	log "github.com/Nixonxp/discord/auth/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"math/rand"
	"time"
)

type Deps struct {
	// deprecated
	UserRepo    UsersStorage
	UserService UsecaseServiceInterface
	Log         *log.Logger
	Oauth2Cgf   oauth2.Config
}

type AuthUsecase struct {
	Deps
}

var _ UsecaseInterface = (*AuthUsecase)(nil)

func NewAuthUsecase(d Deps) UsecaseInterface {
	return &AuthUsecase{
		Deps: d,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, registerInfo RegisterUserInfo) (*models.User, error) {
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

func (u *AuthUsecase) Login(ctx context.Context, loginInfo LoginUserInfo) (*models.LoginResult, error) {
	user, err := u.UserService.GetUserForLogin(ctx, loginInfo)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
	if err != nil {
		return nil, models.ErrCredInvalid
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
	token, err := GenerateJWTToken(userId)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshJWTToken(userId)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func (u *AuthUsecase) Refresh(_ context.Context, refreshToken string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshSecretKey), nil
	})
	if err != nil {
		return "", pkgerrors.Wrap("Invalid or expired JWT token", models.Unauthenticated)
	}

	if !token.Valid {
		return "", pkgerrors.Wrap("Invalid JWT token", models.Unauthenticated)
	}

	claims := token.Claims.(*jwt.StandardClaims)
	userID := claims.Id

	newToken, err := GenerateJWTToken(userID)
	if !token.Valid {
		return "", pkgerrors.Wrap("error generate JWT token", models.Unauthenticated)
	}

	return newToken, nil
}

const secretKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"               // todo to env
const refreshSecretKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9Refresh" // todo to env

func GenerateJWTToken(userID string) (string, error) {
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

func GenerateRefreshJWTToken(userID string) (string, error) {
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

func (u *AuthUsecase) OauthLogin(_ context.Context, _ OauthLoginRequest) (*models.OauthLoginResult, error) {
	url := u.Oauth2Cgf.AuthCodeURL("state")
	return &models.OauthLoginResult{
		Code: url,
	}, nil
}

type UserInfo struct {
	OauthId        string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	AvatarPhotoUrl string `json:"picture"`
}

func generateRandomString(length int) string { // todo to pkg
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}

func (u *AuthUsecase) OauthLoginCallback(ctx context.Context, req OauthLoginCallbackRequest) (*models.LoginResult, error) {
	if req.State != "state" {
		return nil, pkgerrors.Wrap("Invalid state parameter", models.Unauthenticated)
	}

	token, err := u.Oauth2Cgf.Exchange(context.Background(), req.Code)
	if err != nil {
		u.Log.WithContext(ctx).WithError(err).WithField("token", token).Error("Failed to exchange token")
		return nil, pkgerrors.Wrap("Failed to exchange token", models.Unauthenticated)
	}

	client := u.Oauth2Cgf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		u.Log.WithContext(ctx).WithError(err).WithField("token", token).Error("fail auth client from token")
		return nil, pkgerrors.Wrap("Failed to get user info", models.Unauthenticated)
	}

	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	var userInfo UserInfo
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		u.Log.WithContext(ctx).WithError(err).WithField("token", token).Error("Failed unpack user info")
		return nil, pkgerrors.Wrap("Failed unpack user info", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(generateRandomString(10)), bcrypt.MinCost)
	if err != nil {
		return nil, pkgerrors.Wrap("password hashing error", err)
	}

	user, err := u.UserService.CreateOrCreateUser(ctx, GetOrCreateUserRequest{
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
