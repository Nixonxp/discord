package user

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

func (s *UserClient) Login(ctx context.Context, loginInfo usecases.LoginUserInfo) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user_service.Login")
	defer span.Finish()
	response, err := s.client.GetUserByLoginAndPassword(ctx, &user.GetUserByLoginAndPasswordRequest{
		Login:    loginInfo.Login,
		Password: loginInfo.Password,
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
