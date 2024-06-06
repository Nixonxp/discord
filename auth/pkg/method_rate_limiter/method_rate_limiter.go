package method_rate_limiter

import (
	"context"
	"github.com/Nixonxp/discord/auth/pkg/rate_limiter"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type MethodLimiterInfo struct {
	methodName string
	limiter    *rate_limiter.RateLimiter
}

func NewMethodLimiterInfo(methodName string, limitPerMinute int) *MethodLimiterInfo {
	return &MethodLimiterInfo{
		methodName: methodName,
		limiter:    rate_limiter.NewRateLimiter(limitPerMinute),
	}
}

type MethodRateLimiterInterceptor struct {
	limiters map[string]*rate_limiter.RateLimiter
}

func NewMethodRateLimiterInterceptor(limiters ...*MethodLimiterInfo) *MethodRateLimiterInterceptor {
	limiterSet := make(map[string]*rate_limiter.RateLimiter, len(limiters))
	for _, v := range limiters {
		limiterSet[v.methodName] = v.limiter
	}
	return &MethodRateLimiterInterceptor{
		limiters: limiterSet,
	}
}

func (mrli *MethodRateLimiterInterceptor) GetInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if limiter, exist := mrli.checkLimiter(info.FullMethod); exist {
			if limiter.Limit() {
				return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc_ratelimit middleware, please retry later.", info.FullMethod)
			}
		}

		return handler(ctx, req)
	}
}

func (mrli *MethodRateLimiterInterceptor) checkLimiter(method string) (ratelimit.Limiter, bool) {
	for methodName, limiter := range mrli.limiters {
		if strings.Contains(strings.ToLower(method), methodName) {
			return limiter, true
		}
	}

	return nil, false
}
