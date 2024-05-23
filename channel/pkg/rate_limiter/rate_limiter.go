package rate_limiter

import (
	"fmt"
	"github.com/juju/ratelimit"
	"time"
)

type RateLimiter struct {
	TokenBucket *ratelimit.Bucket
}

func NewRateLimiter(tokenCount int) *RateLimiter {
	return &RateLimiter{
		TokenBucket: ratelimit.NewBucket(time.Minute, int64(tokenCount)),
	}
}

func (r *RateLimiter) Limit() bool {
	fmt.Printf("Token Avail %d \n", r.TokenBucket.Available())

	tokenRes := r.TokenBucket.TakeAvailable(1)
	if tokenRes == 0 {
		fmt.Printf("Reached Rate-Limiting %d \n", r.TokenBucket.Available())
		return true
	}

	return false
}
