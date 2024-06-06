package middleware

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/user/internal/app/models"
	log "github.com/Nixonxp/discord/user/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorsUnaryInterceptor - convert any arror to rpc error
func ErrorsUnaryInterceptor(log *log.Logger) grpc.UnaryServerInterceptor {
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
			log.WithContext(ctx).WithError(err).Warn("already exists")
			err = status.Error(codes.AlreadyExists, err.Error())
		case errors.Is(err, models.ErrUnimplemented):
			log.WithContext(ctx).WithError(err).Error("method not implemented")
			err = status.Error(codes.Unimplemented, err.Error())
		case errors.Is(err, models.ErrNotFound):
			log.WithContext(ctx).WithError(err).Warn("not found")
			err = status.Error(codes.NotFound, err.Error())
		case errors.Is(err, models.ErrCredInvalid):
			log.WithContext(ctx).WithError(err).Warn("error invalid credentials")
			err = status.Error(codes.Unauthenticated, err.Error())
		case errors.Is(err, models.Unauthenticated):
			log.WithContext(ctx).WithError(err).Warn("error invalid credentials")
			err = status.Error(codes.Unauthenticated, err.Error())
		case errors.Is(err, models.PermissionDenied):
			log.WithContext(ctx).WithError(err).Warn("error permission denied")
			err = status.Error(codes.PermissionDenied, err.Error())
		default:
			log.WithContext(ctx).WithError(err).Error("internal error")
			err = status.Error(codes.Internal, err.Error())
		}

		return
	}
}
