package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	sync.Mutex
	requests map[string][]time.Time
	window   time.Duration
	limit    int
}

var limiter = &rateLimiter{
	requests: make(map[string][]time.Time),
	window:   time.Minute, // 1 minute window
	limit:    100,         // 100 requests per minute
}

// RateLimiter limits the number of requests per IP
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter.Lock()
		defer limiter.Unlock()

		now := time.Now()
		windowStart := now.Add(-limiter.window)

		// Remove old requests
		var recent []time.Time
		for _, t := range limiter.requests[ip] {
			if t.After(windowStart) {
				recent = append(recent, t)
			}
		}
		limiter.requests[ip] = recent

		// Check if limit exceeded
		if len(recent) >= limiter.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"retry_after": time.Until(recent[0].Add(limiter.window)).
					Round(time.Second).String(),
			})
			c.Abort()
			return
		}

		// Add current request
		limiter.requests[ip] = append(limiter.requests[ip], now)
		c.Next()
	}
}
