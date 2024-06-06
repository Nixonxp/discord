package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
)

func GetUserIdFromContext(ctx context.Context) (string, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	_ = md

	userData := md.Get("userId")
	if userData == nil {
		return "", errors.New("user data not found")
	}

	if len(userData) == 0 {
		return "", errors.New("user data not found")
	}

	userId := userData[0]

	if userId == "" {
		return "", errors.New("user data not found")
	}

	return userId, nil
}
