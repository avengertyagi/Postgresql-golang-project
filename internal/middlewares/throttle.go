package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type throttleState struct {
	count       int
	windowStart time.Time
	blockedAt   time.Time
}

var (
	throttleStatesMu sync.Mutex
	throttleStates   = make(map[string]*throttleState)
	throttleOnce     sync.Once
)

func startCleanupTask() {
	throttleOnce.Do(func() {
		go func() {
			for {
				time.Sleep(5 * time.Minute)
				throttleStatesMu.Lock()
				now := time.Now()
				for key, state := range throttleStates {
					if !state.blockedAt.IsZero() {
						if now.Sub(state.blockedAt) > 5*time.Minute {
							delete(throttleStates, key)
						}
					} else if !state.windowStart.IsZero() && now.Sub(state.windowStart) > 5*time.Minute {
						delete(throttleStates, key)
					}
				}
				throttleStatesMu.Unlock()
			}
		}()
	})
}

func Throttle(maxAttempts int, decayMinutes int) gin.HandlerFunc {
	startCleanupTask()
	decayDuration := time.Duration(decayMinutes) * time.Minute
	return func(c *gin.Context) {
		key := fmt.Sprintf("all:%s:%s", c.ClientIP(), c.Request.URL.Path)
		throttleStatesMu.Lock()
		state, exists := throttleStates[key]
		if !exists {
			state = &throttleState{}
			throttleStates[key] = state
		}
		now := time.Now()
		if !state.blockedAt.IsZero() {
			if now.Sub(state.blockedAt) < decayDuration {
				throttleStatesMu.Unlock()
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"status":     false,
					"statusCode": http.StatusTooManyRequests,
					"message":    fmt.Sprintf("Too many requests. Please try again after %d minutes.", decayMinutes),
				})
				return
			}
			state.count = 0
			state.blockedAt = time.Time{}
			state.windowStart = time.Time{}
		}
		if state.windowStart.IsZero() || now.Sub(state.windowStart) >= decayDuration {
			state.windowStart = now
			state.count = 0
		}
		state.count++
		if state.count >= maxAttempts {
			state.blockedAt = now
		}
		throttleStatesMu.Unlock()
		c.Next()
	}
}

func ThrottleFailures(maxAttempts int, decayMinutes int) gin.HandlerFunc {
	startCleanupTask()
	decayDuration := time.Duration(decayMinutes) * time.Minute
	return func(c *gin.Context) {
		key := fmt.Sprintf("fail:%s:%s", c.ClientIP(), c.Request.URL.Path)
		throttleStatesMu.Lock()
		state, exists := throttleStates[key]
		if exists && !state.blockedAt.IsZero() {
			if time.Since(state.blockedAt) < decayDuration {
				throttleStatesMu.Unlock()
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"status":     false,
					"statusCode": http.StatusTooManyRequests,
					"message":    fmt.Sprintf("Too many failed attempts. Please try again after %d minutes.", decayMinutes),
				})
				return
			}
			state.count = 0
			state.blockedAt = time.Time{}
		}
		throttleStatesMu.Unlock()
		c.Next()
		status := c.Writer.Status()
		if status == http.StatusUnauthorized {
			throttleStatesMu.Lock()
			state, exists := throttleStates[key]
			if !exists {
				state = &throttleState{}
				throttleStates[key] = state
			}
			state.count++
			if state.count >= maxAttempts {
				state.blockedAt = time.Now()
			}
			throttleStatesMu.Unlock()
		} else if status == http.StatusOK {
			throttleStatesMu.Lock()
			delete(throttleStates, key)
			throttleStatesMu.Unlock()
		}
	}
}
