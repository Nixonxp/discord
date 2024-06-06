package middleware

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/auth/internal/app/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorsUnaryInterceptor - convert any error to rpc error
func ErrorsUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		if _, ok := status.FromError(err); ok {
			return
		}

		switch {
		case errors.Is(err, models.ErrAlreadyExists):
			err = status.Error(codes.AlreadyExists, err.Error())
		case errors.Is(err, models.ErrUnimplemented):
			err = status.Error(codes.Unimplemented, err.Error())
		case errors.Is(err, models.ErrCredInvalid):
			err = status.Error(codes.Unauthenticated, err.Error())
		case errors.Is(err, models.Unauthenticated):
			//log.WithContext(ctx).WithError(err).Warn("error invalid credentials") todo add logger
			err = status.Error(codes.Unauthenticated, err.Error())
		default:
			err = status.Error(codes.Internal, err.Error())
		}

		return
	}
}
