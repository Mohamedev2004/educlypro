package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"educlypro/shared/types"
	"educlypro/shared/utils"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"gorm.io/gorm"
)

type Service interface {
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error)
	Logout(ctx context.Context, userID uint) error
	Me(ctx context.Context, userID uint) (*UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest) (*UserResponse, error)
	UpdatePassword(ctx context.Context, userID uint, req *UpdatePasswordRequest) error
	ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
}

type service struct {
	repo      Repository
	publisher message.Publisher
}

func NewService(db *gorm.DB, publisher message.Publisher) Service {
	repo := NewRepository(db)
	return &service{repo: repo, publisher: publisher}
}

func (s *service) UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest) (*UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "User", err)
		return nil, err
	}
	if user == nil {
		s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "User", fmt.Sprint(userID), "system", http.StatusInternalServerError, map[string]string{
			"user_id": fmt.Sprint(userID),
			"error":   "valid token but user missing from db",
		}, "system.events.v1.auth.anomaly_user_missing")
		return nil, errors.New("user not found")
	}

	if req.Email != user.Email {
		existing, err := s.repo.FindByEmail(req.Email)
		if err != nil {
			s.publishError(ctx, types.ActionFailed, "User", err)
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("email already in use")
		}
	}

	if err := s.repo.UpdateProfile(userID, req.Username, req.Email); err != nil {
		s.publishError(ctx, types.ActionFailed, "User", err)
		return nil, err
	}

	s.publishEvent(ctx, types.LevelInfo, types.ActionUpdated, "User", fmt.Sprint(user.ID), fmt.Sprint(user.ID), http.StatusOK, map[string]string{
		"username": req.Username,
		"email":    req.Email,
	}, "system.events.v1.auth.profile_updated")

	return &UserResponse{
		ID:       user.ID,
		Username: req.Username,
		Email:    req.Email,
		Role:     user.Role.Name,
		Center:   centerResponse(user),
	}, nil
}

func (s *service) UpdatePassword(ctx context.Context, userID uint, req *UpdatePasswordRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "User", err)
		return err
	}
	if user == nil {
		s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "User", fmt.Sprint(userID), "system", http.StatusInternalServerError, map[string]string{
			"user_id": fmt.Sprint(userID),
			"error":   "valid token but user missing from db",
		}, "system.events.v1.auth.anomaly_user_missing")
		return errors.New("user not found")
	}

	if !utils.CheckPassword(req.CurrentPassword, user.Password) {
		s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "User", fmt.Sprint(user.ID), fmt.Sprint(user.ID), http.StatusBadRequest, map[string]string{
			"username": user.Username,
			"email":    user.Email,
			"reason":   "current password mismatch",
		}, "system.events.v1.auth.password_update_failed")
		return errors.New("invalid current password")
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "User", err)
		return err
	}

	if err := s.repo.UpdatePassword(userID, hashed); err != nil {
		s.publishError(ctx, types.ActionFailed, "User", err)
		return err
	}

	s.publishEvent(ctx, types.LevelInfo, types.ActionUpdated, "User", fmt.Sprint(user.ID), fmt.Sprint(user.ID), http.StatusOK, map[string]string{
		"username": user.Username,
		"email":    user.Email,
	}, "system.events.v1.auth.password_updated")

	return nil
}

func (s *service) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "User", err)
		return err
	}

	if user == nil {
		s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "AuthSession", "unknown", "unknown", http.StatusOK, map[string]string{
			"attempted_email": req.Email,
			"reason":          "email not found in db",
		}, "system.events.v1.auth.password_reset_unknown_email")
		return nil // Silent fail for security
	}

	tokenStr := utils.GenerateUUID()
	hashedToken := utils.HashSHA256(tokenStr)

	token := &Token{
		Token:     hashedToken,
		UserID:    user.ID,
		Type:      "reset_password",
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	if err := s.repo.CreateToken(token); err != nil {
		s.publishError(ctx, types.ActionFailed, "Token", err)
		return err
	}

	// NOTE: token is included here to trigger the email handler — it is NOT stored in audit logs in plaintext
	s.publishEvent(ctx, types.LevelInfo, types.ActionCreated, "User", fmt.Sprint(user.ID), fmt.Sprint(user.ID), http.StatusOK, map[string]string{
		"username":   user.Username,
		"email":      user.Email,
		"token":      tokenStr,
		"expires_in": "30 minutes",
	}, "system.events.v1.auth.password_reset_requested")

	return nil
}

func (s *service) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	hashedToken := utils.HashSHA256(req.Token)
	token, err := s.repo.FindByToken(hashedToken)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "Token", err)
		return err
	}

	if token == nil || token.Type != "reset_password" || token.ExpiresAt.Before(time.Now()) {
		s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "User", "unknown", "unknown", http.StatusBadRequest, map[string]string{
			"reason": "token is nil, wrong type, or expired",
		}, "system.events.v1.auth.password_reset_failed")
		return errors.New("invalid or expired token")
	}

	user, err := s.repo.FindByID(token.UserID)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "User", err)
		return err
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "User", err)
		return err
	}

	if err := s.repo.UpdatePassword(token.UserID, hashed); err != nil {
		s.publishError(ctx, types.ActionFailed, "User", err)
		return err
	}

	if err := s.repo.DeleteTokensByUserID(token.UserID, "reset_password"); err != nil {
		s.publishError(ctx, types.ActionFailed, "Token", err)
		return err
	}

	payload := map[string]string{
		"user_id": fmt.Sprint(token.UserID),
	}
	if user != nil {
		payload["username"] = user.Username
		payload["email"] = user.Email
	}

	s.publishEvent(ctx, types.LevelInfo, types.ActionUpdated, "User", fmt.Sprint(token.UserID), fmt.Sprint(token.UserID), http.StatusOK, payload, "system.events.v1.auth.password_reset_completed")

	return nil
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "AuthSession", err)
		return nil, err
	}
	if user == nil {
		s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "AuthSession", "unknown", "unknown", http.StatusUnauthorized, map[string]string{
			"attempted_email": req.Email,
			"reason":          "email not found",
		}, "system.events.v1.auth.login_failed")
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "AuthSession", fmt.Sprint(user.ID), "unknown", http.StatusUnauthorized, map[string]string{
			"attempted_email": req.Email,
			"username":        user.Username,
			"reason":          "wrong password",
		}, "system.events.v1.auth.login_failed")
		return nil, errors.New("invalid credentials")
	}

	_ = s.repo.DeleteTokensByUserID(user.ID, "access")
	_ = s.repo.DeleteTokensByUserID(user.ID, "refresh")

	tokens, err := s.issueAuthTokens(user, user.Role.Name)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "Token", err)
		return nil, err
	}

	s.publishEvent(ctx, types.LevelInfo, types.ActionUpdated, "AuthSession", fmt.Sprint(user.ID), fmt.Sprint(user.ID), http.StatusOK, map[string]string{
		"username": user.Username,
		"email":    user.Email,
	}, "system.events.v1.auth.logged_in")

	return &AuthResponse{
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role.Name,
			Center:   centerResponse(user),
		},
		Token:        tokens.Token,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token required")
	}

	claims, err := utils.ParseJWT(refreshToken)
	if err != nil || claims.TokenType != "refresh" {
		s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "AuthSession", "unknown", "unknown", http.StatusUnauthorized, map[string]string{
			"error":  "malformed or invalid jwt",
			"reason": "token parse failed or wrong type",
		}, "system.events.v1.auth.refresh_failed")
		return nil, errors.New("invalid refresh token")
	}

	hashedToken := utils.HashSHA256(refreshToken)
	tokenRecord, err := s.repo.FindByToken(hashedToken)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "Token", err)
		return nil, err
	}
	if tokenRecord == nil || tokenRecord.Type != "refresh" {
		s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "AuthSession", fmt.Sprint(claims.UserID), "unknown", http.StatusUnauthorized, map[string]string{
			"user_id": fmt.Sprint(claims.UserID),
			"error":   "token revoked or not in db",
		}, "system.events.v1.auth.refresh_failed")
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.repo.FindByID(claims.UserID)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "User", err)
		return nil, err
	}
	if user == nil {
		s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "User", fmt.Sprint(claims.UserID), "system", http.StatusInternalServerError, map[string]string{
			"user_id": fmt.Sprint(claims.UserID),
			"error":   "valid token but user missing from db",
		}, "system.events.v1.auth.anomaly_user_missing")
		return nil, errors.New("user not found")
	}

	if err := s.repo.DeleteToken(hashedToken, "refresh"); err != nil {
		s.publishError(ctx, types.ActionFailed, "Token", err)
		return nil, err
	}
	_ = s.repo.DeleteTokensByUserID(user.ID, "access")

	tokens, err := s.issueAuthTokens(user, user.Role.Name)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "Token", err)
		return nil, err
	}

	s.publishEvent(ctx, types.LevelInfo, types.ActionUpdated, "AuthSession", fmt.Sprint(user.ID), fmt.Sprint(user.ID), http.StatusOK, map[string]string{
		"username": user.Username,
		"email":    user.Email,
	}, "system.events.v1.auth.token_refreshed")

	return &AuthResponse{
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role.Name,
			Center:   centerResponse(user),
		},
		Token:        tokens.Token,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *service) Logout(ctx context.Context, userID uint) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "Token", err)
		return err
	}

	if err := s.repo.DeleteTokensByUserID(userID, "access"); err != nil {
		s.publishError(ctx, types.ActionFailed, "Token", err)
		return err
	}
	if err := s.repo.DeleteTokensByUserID(userID, "refresh"); err != nil {
		s.publishError(ctx, types.ActionFailed, "Token", err)
		return err
	}

	payload := map[string]string{
		"user_id": fmt.Sprint(userID),
	}
	if user != nil {
		payload["username"] = user.Username
		payload["email"] = user.Email
	}

	s.publishEvent(ctx, types.LevelInfo, types.ActionUpdated, "AuthSession", fmt.Sprint(userID), fmt.Sprint(userID), http.StatusOK, payload, "system.events.v1.auth.logged_out")

	return nil
}

func (s *service) Me(ctx context.Context, userID uint) (*UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, "User", err)
		return nil, err
	}
	if user == nil {
		s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "User", fmt.Sprint(userID), "system", http.StatusNotFound, map[string]string{
			"user_id": fmt.Sprint(userID),
			"error":   "valid token but user missing from db",
		}, "system.events.v1.auth.anomaly_user_missing")
		return nil, errors.New("user not found")
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role.Name,
		Center:   centerResponse(user),
	}, nil
}

func centerResponse(user *User) *CenterResponse {
	if user.Center == nil {
		return nil
	}
	return &CenterResponse{ID: user.Center.ID, Name: user.Center.Name, Slug: user.Center.Slug}
}

func (s *service) issueAuthTokens(user *User, role string) (*AuthResponse, error) {
	accessExpiry := utils.GetAccessTokenExpiry()
	refreshExpiry := utils.GetRefreshTokenExpiry()

	accessToken, err := utils.GenerateJWT(user.ID, role, "access", accessExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateJWT(user.ID, role, "refresh", refreshExpiry)
	if err != nil {
		return nil, err
	}

	hashedAccess := utils.HashSHA256(accessToken)
	hashedRefresh := utils.HashSHA256(refreshToken)

	if err := s.repo.CreateToken(&Token{
		Token:     hashedAccess,
		UserID:    user.ID,
		Type:      "access",
		ExpiresAt: time.Now().Add(accessExpiry),
	}); err != nil {
		return nil, err
	}

	if err := s.repo.CreateToken(&Token{
		Token:     hashedRefresh,
		UserID:    user.ID,
		Type:      "refresh",
		ExpiresAt: time.Now().Add(refreshExpiry),
	}); err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ==========================================
// INTERNAL HELPERS FOR PUBLISHING EVENTS
// ==========================================

func (s *service) publishError(ctx context.Context, action types.EventAction, entity string, err error) {
	if err == nil {
		return
	}
	payload := map[string]string{
		"error_message": err.Error(),
	}
	s.publishEvent(ctx, types.LevelError, action, entity, "unknown", "system", http.StatusInternalServerError, payload, "system.events.v1.auth.system_error")
}

func (s *service) publishEvent(
	ctx context.Context,
	level types.LogLevel,
	action types.EventAction,
	entity string,
	entityID string,
	actorID string,
	statusCode int,
	payload any,
	topic string,
) {
	if s.publisher == nil {
		return
	}

	reqID, _ := ctx.Value("X-Request-ID").(string)

	start, _ := ctx.Value("startTime").(time.Time)

	var duration float64 = 0
	if !start.IsZero() {
		duration = math.Round(time.Since(start).Seconds()*100000) / 100
	}

	event := types.AuditEvent{
		RequestID:  reqID,
		Level:      level,
		Entity:     entity,
		EntityID:   entityID,
		ActorID:    actorID,
		StatusCode: statusCode,
		Action:     action,
		Payload:    payload,
		Timestamp:  time.Now(),
		DurationMs: duration,
	}

	payloadBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event payload: %v", err)
		return
	}

	specificMsg := message.NewMessage(watermill.NewUUID(), payloadBytes)
	if err := s.publisher.Publish(topic, specificMsg); err != nil {
		log.Printf("Failed to publish event to topic %s: %v", topic, err)
	}

	logMsg := message.NewMessage(watermill.NewUUID(), payloadBytes)
	if err := s.publisher.Publish("system.audit_logs", logMsg); err != nil {
		log.Printf("Failed to publish to audit logs: %v", err)
	}
}
