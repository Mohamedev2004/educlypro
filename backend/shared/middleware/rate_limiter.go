package middleware

import (
	"encoding/json"
	"math"
	"net/http"
	"sync"
	"time"

	"educlypro/shared/types"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex

	once sync.Once
)

// cleanup old IPs to prevent memory leak
func cleanupClients() {
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
}

// RateLimiterMiddleware creates a middleware with custom rate + burst
func RateLimiterMiddleware(rps rate.Limit, burst int, publisher message.Publisher) gin.HandlerFunc {

	// ensure cleanup runs only once
	once.Do(func() {
		go cleanupClients()
	})

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()

		if _, exists := clients[ip]; !exists {
			clients[ip] = &client{
				limiter: rate.NewLimiter(rps, burst),
			}
		}

		clients[ip].lastSeen = time.Now()
		limiter := clients[ip].limiter

		mu.Unlock()

		if !limiter.Allow() {
			// Publish an audit event (best-effort) when a request is blocked.
			if publisher != nil {
				reqID, _ := c.Request.Context().Value("X-Request-ID").(string)
				start, _ := c.Request.Context().Value("startTime").(time.Time)

				var duration float64 = 0
				if !start.IsZero() {
					duration = math.Round(time.Since(start).Seconds()*100000) / 100
				}

				event := types.AuditEvent{
					RequestID:  reqID,
					Level:      types.LevelWarn,
					Entity:     "RateLimiter",
					EntityID:   ip,
					ActorID:    "unknown",
					StatusCode: http.StatusTooManyRequests,
					Action:     types.ActionFailed,
					Payload:    gin.H{"ip": ip, "method": c.Request.Method, "path": c.FullPath(), "rps": float64(rps), "burst": burst},
					Timestamp:  time.Now(),
					DurationMs: duration,
				}

				if payloadBytes, err := json.Marshal(event); err == nil {
					logMsg := message.NewMessage(watermill.NewUUID(), payloadBytes)
					_ = publisher.Publish("system.audit_logs", logMsg)
				}
			}

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please slow down.",
			})
			return
		}

		c.Next()
	}
}
