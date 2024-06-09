package services

import (
	"context"
	"github.com/Nixonxp/discord/auth/internal/app/models"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	"github.com/Nixonxp/discord/auth/pkg/api/user"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

func (s *UserClient) Register(ctx context.Context, registerInfo usecases.RegisterUserInfo) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user_service.Register")
	defer span.Finish()
	response, err := s.client.CreateUser(ctx, &user.CreateUserRequest{
		Login:    registerInfo.Login,
		Name:     registerInfo.Name,
		Email:    registerInfo.Email,
		Password: registerInfo.Password,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		UserID: models.UserID(uuid.MustParse(response.Id)),
		Login:  response.Login,
		Name:   response.Name,
		Email:  response.Email,
	}, nil
}

func (s *UserClient) GetUserForLogin(ctx context.Context, loginInfo usecases.LoginUserInfo) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user_service.GetUserForLogin")
	defer span.Finish()
	response, err := s.client.GetUserForLogin(ctx, &user.GetUserForLoginRequest{
		Login: loginInfo.Login,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		UserID:   models.UserID(uuid.MustParse(response.Id)),
		Login:    response.Login,
		Name:     response.Name,
		Email:    response.Email,
		Password: response.Password,
	}, nil
}

func (s *UserClient) CreateOrCreateUser(ctx context.Context, userInfo usecases.GetOrCreateUserRequest) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user_service.CreateOrCreateUser")
	defer span.Finish()
	response, err := s.client.CreateOrGetUser(ctx, &user.CreateOrGetUserRequest{
		Login:          userInfo.Login,
		Name:           userInfo.Name,
		Email:          userInfo.Email,
		Password:       userInfo.Password,
		AvatarPhotoUrl: userInfo.AvatarPhotoUrl,
		OauthId:        userInfo.OauthId,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		UserID:         models.UserID(uuid.MustParse(response.Id)),
		Login:          response.Login,
		Name:           response.Name,
		Email:          response.Email,
		AvatarPhotoUrl: response.AvatarPhotoUrl,
	}, nil
}
