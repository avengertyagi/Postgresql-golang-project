package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type rateLimitClient struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func RateLimiter() gin.HandlerFunc {
	var (
		mu      sync.Mutex
		clients = make(map[string]*rateLimitClient)
	)
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		if _, exists := clients[ip]; !exists {
			clients[ip] = &rateLimitClient{
				limiter: rate.NewLimiter(10, 20),
			}
		}
		cl := clients[ip]
		cl.lastSeen = time.Now()
		mu.Unlock()
		if !cl.limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":     false,
				"statusCode": http.StatusTooManyRequests,
				"error":      "Too many requests. Please try again later.",
			})
			return
		}
		c.Next()
	}
}
