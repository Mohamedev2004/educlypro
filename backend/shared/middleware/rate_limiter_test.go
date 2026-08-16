package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockPublisher struct {
	published []string
}

func (m *mockPublisher) Publish(topic string, messages ...*message.Message) error {
	m.published = append(m.published, topic)
	return nil
}

func (m *mockPublisher) Close() error { return nil }

// ── Helpers ───────────────────────────────────────────────────────────────────

// resetClients clears the package-level client map so tests don't leak
// limiter state into one another via shared IPs.
func resetClients() {
	mu.Lock()
	clients = make(map[string]*client)
	mu.Unlock()
}

func setupRouter(rps rate.Limit, burst int, publisher message.Publisher) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimiterMiddleware(rps, burst, publisher))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return r
}

func doRequest(r *gin.Engine, ip string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = ip + ":12345"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestRateLimiterMiddleware_AllowsWithinBurst(t *testing.T) {
	resetClients()
	r := setupRouter(rate.Limit(1), 3, nil)

	for i := 1; i <= 3; i++ {
		if w := doRequest(r, "10.0.0.1"); w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, w.Code)
		}
	}
}

func TestRateLimiterMiddleware_BlocksOverBurst(t *testing.T) {
	resetClients()
	r := setupRouter(rate.Limit(1), 3, nil)

	for i := 0; i < 3; i++ {
		doRequest(r, "10.0.0.2")
	}

	w := doRequest(r, "10.0.0.2")
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 once burst is exhausted, got %d", w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if want := "Too many requests. Please slow down."; body["error"] != want {
		t.Fatalf("expected error %q, got %q", want, body["error"])
	}
}

func TestRateLimiterMiddleware_RefillsOverTime(t *testing.T) {
	resetClients()
	// 20 tokens/sec => one token every 50ms, burst of 1.
	r := setupRouter(rate.Limit(20), 1, nil)

	if w := doRequest(r, "10.0.0.3"); w.Code != http.StatusOK {
		t.Fatalf("expected first request to be allowed, got %d", w.Code)
	}
	if w := doRequest(r, "10.0.0.3"); w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected immediate follow-up to be blocked, got %d", w.Code)
	}

	time.Sleep(100 * time.Millisecond)

	if w := doRequest(r, "10.0.0.3"); w.Code != http.StatusOK {
		t.Fatalf("expected request after refill window to be allowed, got %d", w.Code)
	}
}

func TestRateLimiterMiddleware_TracksClientsIndependently(t *testing.T) {
	resetClients()
	r := setupRouter(rate.Limit(1), 1, nil)

	doRequest(r, "10.0.0.4")
	if w := doRequest(r, "10.0.0.4"); w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected client A to be rate-limited, got %d", w.Code)
	}

	if w := doRequest(r, "10.0.0.5"); w.Code != http.StatusOK {
		t.Fatalf("expected client B to be unaffected by client A's limit, got %d", w.Code)
	}
}

func TestRateLimiterMiddleware_PublishesAuditEventOnBlock(t *testing.T) {
	resetClients()
	pub := &mockPublisher{}
	r := setupRouter(rate.Limit(1), 1, pub)

	doRequest(r, "10.0.0.6") // consumes the only token
	doRequest(r, "10.0.0.6") // blocked, should publish an audit event

	if len(pub.published) != 1 {
		t.Fatalf("expected exactly 1 published event, got %d", len(pub.published))
	}
	if want := "system.audit_logs"; pub.published[0] != want {
		t.Fatalf("expected event published to %q, got %q", want, pub.published[0])
	}
}

func TestRateLimiterMiddleware_NoPublisherDoesNotPanic(t *testing.T) {
	resetClients()
	r := setupRouter(rate.Limit(1), 1, nil)

	doRequest(r, "10.0.0.7")
	if w := doRequest(r, "10.0.0.7"); w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 even without a publisher, got %d", w.Code)
	}
}
