package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"educlypro/config"
	"educlypro/shared/utils"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
)

func init() {
	os.Setenv("JWT_SECRET", "test-secret-key-at-least-32-chars-long-123456")
	config.Load()
	hash, _ := utils.HashPassword("secret")
	existingUser.Password = hash
}

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockPublisher struct {
	publishedTopics []string
}

func (m *mockPublisher) Publish(topic string, messages ...*message.Message) error {
	m.publishedTopics = append(m.publishedTopics, topic)
	return nil
}

func (m *mockPublisher) Close() error { return nil }

type mockService struct {
	loginFn          func(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	refreshFn        func(ctx context.Context, refreshToken string) (*AuthResponse, error)
	logoutFn         func(ctx context.Context, userID uint) error
	meFn             func(ctx context.Context, userID uint) (*UserResponse, error)
	updatePasswordFn func(ctx context.Context, userID uint, req *UpdatePasswordRequest) error
	updateProfileFn  func(ctx context.Context, userID uint, req *UpdateProfileRequest) (*UserResponse, error)
	forgotPasswordFn func(ctx context.Context, req *ForgotPasswordRequest) error
	resetPasswordFn  func(ctx context.Context, req *ResetPasswordRequest) error
}

func (m *mockService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	return m.loginFn(ctx, req)
}
func (m *mockService) Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	return m.refreshFn(ctx, refreshToken)
}
func (m *mockService) Logout(ctx context.Context, userID uint) error { return m.logoutFn(ctx, userID) }
func (m *mockService) Me(ctx context.Context, userID uint) (*UserResponse, error) {
	return m.meFn(ctx, userID)
}
func (m *mockService) UpdatePassword(ctx context.Context, userID uint, req *UpdatePasswordRequest) error {
	return m.updatePasswordFn(ctx, userID, req)
}
func (m *mockService) UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest) (*UserResponse, error) {
	return m.updateProfileFn(ctx, userID, req)
}
func (m *mockService) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error {
	return m.forgotPasswordFn(ctx, req)
}
func (m *mockService) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	return m.resetPasswordFn(ctx, req)
}

type mockRepository struct {
	findByEmailFn        func(email string) (*User, error)
	findByIDFn           func(id uint) (*User, error)
	createFn             func(user *User) error
	createTokenFn        func(token *Token) error
	deleteTokensByUserFn func(userID uint, tokenType string) error
	deleteTokenFn        func(tokenStr string, tokenType string) error
	findByTokenFn        func(tokenStr string) (*Token, error)
}

func (m *mockRepository) FindByEmail(email string) (*User, error) { return m.findByEmailFn(email) }
func (m *mockRepository) FindByID(id uint) (*User, error)         { return m.findByIDFn(id) }
func (m *mockRepository) Create(user *User) error                 { return m.createFn(user) }
func (m *mockRepository) CreateToken(token *Token) error {
	if m.createTokenFn != nil {
		return m.createTokenFn(token)
	}
	return nil
}
func (m *mockRepository) DeleteTokensByUserID(userID uint, tokenType string) error {
	if m.deleteTokensByUserFn != nil {
		return m.deleteTokensByUserFn(userID, tokenType)
	}
	return nil
}
func (m *mockRepository) DeleteToken(tokenStr string, tokenType string) error {
	return m.deleteTokenFn(tokenStr, tokenType)
}
func (m *mockRepository) FindByToken(tokenStr string) (*Token, error) {
	return m.findByTokenFn(tokenStr)
}
func (m *mockRepository) UpdatePassword(userID uint, hashedPassword string) error        { return nil }
func (m *mockRepository) UpdateProfile(userID uint, username string, email string) error { return nil }

// ── Fixtures & Helpers ────────────────────────────────────────────────────────

var fakeAuthResp = &AuthResponse{
	User:  UserResponse{ID: 1, Username: "john", Email: "john@example.com", Role: "center_owner"},
	Token: "fake.jwt.token", RefreshToken: "fake.refresh.token",
}

var existingUser = &User{
	ID: 1, Username: "john", Email: "john@example.com",
	Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
	Role:     Role{ID: 1, Name: "center_owner"},
}

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout", func(c *gin.Context) { c.Set("userID", uint(1)); h.Logout(c) })
	r.POST("/auth/forgot-password", h.ForgotPassword)
	return r
}

func makeJSON(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

func doRequest(r *gin.Engine, method, path string, body *bytes.Buffer) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var req *http.Request
	if body == nil {
		req, _ = http.NewRequest(method, path, nil)
	} else {
		req, _ = http.NewRequest(method, path, body)
	}
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

// ── Handler Tests (TDT) ───────────────────────────────────────────────────────

func TestHandler_Login(t *testing.T) {
	tests := []struct {
		name       string
		req        LoginRequest
		mockResp   *AuthResponse
		mockErr    error
		wantStatus int
	}{
		{"Success", LoginRequest{Email: "j@t.com", Password: "secret12"}, fakeAuthResp, nil, http.StatusOK},
		{"Invalid Credentials", LoginRequest{Email: "j@t.com", Password: "badpass"}, nil, errors.New("invalid credentials"), http.StatusUnauthorized},
		{"Internal Error", LoginRequest{Email: "j@t.com", Password: "secret12"}, nil, errors.New("db down"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				loginFn: func(ctx context.Context, req *LoginRequest) (*AuthResponse, error) { return tt.mockResp, tt.mockErr },
			}
			r := setupRouter(&Handler{service: svc})
			w := doRequest(r, http.MethodPost, "/auth/login", makeJSON(t, tt.req))

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

// ── Service Tests (TDT with Event Publishing Verification) ────────────────────

func TestService_Login(t *testing.T) {
	tests := []struct {
		name      string
		req       *LoginRequest
		mockRepo  *mockRepository
		wantErr   string
		wantEvent string
	}{
		{
			name: "Success",
			req:  &LoginRequest{Email: "john@example.com", Password: "secret"},
			mockRepo: &mockRepository{
				findByEmailFn: func(e string) (*User, error) { return existingUser, nil },
			},
			wantErr:   "",
			wantEvent: "system.events.v1.auth.logged_in",
		},
		{
			name: "Wrong Password",
			req:  &LoginRequest{Email: "john@example.com", Password: "wrong"},
			mockRepo: &mockRepository{
				findByEmailFn: func(e string) (*User, error) { return existingUser, nil },
			},
			wantErr:   "invalid credentials",
			wantEvent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pub := &mockPublisher{}
			svc := &service{repo: tt.mockRepo, publisher: pub}

			_, err := svc.Login(context.Background(), tt.req)

			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Errorf("expected error %q, got %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				// Verify Event Publishing
				found := false
				for _, topic := range pub.publishedTopics {
					if topic == tt.wantEvent {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected event %q, got topics: %v", tt.wantEvent, pub.publishedTopics)
				}
			}
		})
	}
}

func TestService_ForgotPassword_FiresEvent(t *testing.T) {
	pub := &mockPublisher{}
	mockRepo := &mockRepository{
		findByEmailFn: func(email string) (*User, error) { return existingUser, nil },
		createTokenFn: func(token *Token) error { return nil },
	}
	svc := &service{repo: mockRepo, publisher: pub}

	err := svc.ForgotPassword(context.Background(), &ForgotPasswordRequest{Email: "john@example.com"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// CRITICAL EVENT ARCHITECTURE TEST:
	// Verify that requesting a password reset triggers the fan-out event
	found := false
	for _, topic := range pub.publishedTopics {
		if topic == "system.events.v1.auth.password_reset_requested" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected topic 'system.events.v1.auth.password_reset_requested', got topics: %v", pub.publishedTopics)
	}
}
