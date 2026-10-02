package middleware

import (
	"go_dfs/internal/redis"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func RateLimit(
	rateLimiter  *redis.RateLimiter,
	action string,
	limit int,
	window time.Duration,
) func(http.Handler) http.Handler{
	return func (next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			userID,ok:=GetUserID(r)
			if !ok{
				http.Error(w,"invalid user id", http.StatusInternalServerError)
				return
			}

			id,err:=uuid.Parse(userID)
			if err != nil {
				http.Error(w, "invalid user id", http.StatusInternalServerError)
				return
			}

			key:=redis.UserRateLimitKey(action,id)

			allowed,err:=rateLimiter.Allow(
				r.Context(),
				key,
				limit,
				window,
			)
			if err != nil {
				http.Error(w, "rate limiter unavailable", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		}) 
	}
}