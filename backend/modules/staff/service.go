package staff

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"educlypro/modules/auth"
	"educlypro/shared/types"
	"educlypro/shared/utils"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

var (
	ErrCenterMismatch = errors.New("center mismatch")
	ErrEmailTaken     = errors.New("email already in use")
	ErrStaffNotFound  = errors.New("staff not found")
)

// TokenRevoker is the minimal capability the staff service needs from the
// auth module — revoking a deleted staff member's active sessions. Kept as
// a small local interface (rather than depending on the whole auth.Repository)
// so this module only couples to what it actually uses.
type TokenRevoker interface {
	DeleteTokensByUserID(userID uint, tokenType string) error
}

type Service interface {
	List(ctx context.Context, ownerUserID uint, params ListParams) (*ListResponse, error)
	Create(ctx context.Context, ownerUserID uint, req *CreateRequest) (*Response, error)
	Update(ctx context.Context, ownerUserID, staffID uint, req *UpdateRequest) (*Response, error)
	Delete(ctx context.Context, ownerUserID, staffID uint) error
}

type service struct {
	repo         Repository
	tokenRevoker TokenRevoker
	publisher    message.Publisher
}

func NewService(repo Repository, tokenRevoker TokenRevoker, publisher message.Publisher) Service {
	return &service{repo: repo, tokenRevoker: tokenRevoker, publisher: publisher}
}

func (s *service) List(ctx context.Context, ownerUserID uint, params ListParams) (*ListResponse, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	users, err := s.repo.ListByCenter(ctx, centerID, params)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountByCenter(ctx, centerID, params.Role, params.Search)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / params.PerPage
	if int(total)%params.PerPage != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	items := make([]Response, 0, len(users))
	for _, u := range users {
		items = append(items, toResponse(u))
	}

	return &ListResponse{
		Items: items,
		Pagination: PaginationMeta{
			Page:       params.Page,
			PerPage:    params.PerPage,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    params.Page < totalPages,
			HasPrev:    params.Page > 1,
		},
	}, nil
}

func (s *service) Create(ctx context.Context, ownerUserID uint, req *CreateRequest) (*Response, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, err)
		return nil, err
	}
	if req.CenterID != centerID {
		s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "unknown", fmt.Sprint(ownerUserID), http.StatusForbidden, map[string]any{
			"attempted_center_id": req.CenterID,
			"owner_center_id":     centerID,
			"reason":              "center_id in request did not match the owner's own center",
		}, "system.events.v1.staff.center_mismatch")
		return nil, ErrCenterMismatch
	}

	taken, err := s.repo.EmailTaken(ctx, req.Email, 0)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, err)
		return nil, err
	}
	if taken {
		return nil, ErrEmailTaken
	}

	roleID, err := s.repo.FindRoleIDByName(ctx, req.Role)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, err)
		return nil, err
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, err)
		return nil, err
	}

	user := &auth.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
		RoleID:   roleID,
		CenterID: &centerID,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		s.publishError(ctx, types.ActionFailed, err)
		return nil, err
	}
	user.Role = auth.Role{ID: roleID, Name: req.Role}

	s.publishEvent(ctx, types.LevelInfo, types.ActionCreated, fmt.Sprint(user.ID), fmt.Sprint(ownerUserID), http.StatusCreated, map[string]any{
		"username":  user.Username,
		"email":     user.Email,
		"role":      req.Role,
		"center_id": centerID,
	}, "system.events.v1.staff.created")

	resp := toResponse(*user)
	return &resp, nil
}

func (s *service) Update(ctx context.Context, ownerUserID, staffID uint, req *UpdateRequest) (*Response, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, err)
		return nil, err
	}
	if req.CenterID != centerID {
		return nil, ErrCenterMismatch
	}

	existing, err := s.repo.FindByIDInCenter(ctx, centerID, staffID)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, err)
		return nil, err
	}
	if existing == nil {
		return nil, ErrStaffNotFound
	}

	if req.Email != existing.Email {
		taken, err := s.repo.EmailTaken(ctx, req.Email, staffID)
		if err != nil {
			s.publishError(ctx, types.ActionFailed, err)
			return nil, err
		}
		if taken {
			return nil, ErrEmailTaken
		}
	}

	roleID, err := s.repo.FindRoleIDByName(ctx, req.Role)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, err)
		return nil, err
	}

	password := existing.Password
	if req.Password != nil && *req.Password != "" {
		hashed, err := utils.HashPassword(*req.Password)
		if err != nil {
			s.publishError(ctx, types.ActionFailed, err)
			return nil, err
		}
		password = hashed
	}

	updated := &auth.User{
		ID:       staffID,
		Username: req.Username,
		Email:    req.Email,
		Password: password,
		RoleID:   roleID,
	}
	if err := s.repo.Update(ctx, updated); err != nil {
		s.publishError(ctx, types.ActionFailed, err)
		return nil, err
	}
	updated.CenterID = &centerID
	updated.CreatedAt = existing.CreatedAt
	updated.Role = auth.Role{ID: roleID, Name: req.Role}

	s.publishEvent(ctx, types.LevelInfo, types.ActionUpdated, fmt.Sprint(staffID), fmt.Sprint(ownerUserID), http.StatusOK, map[string]any{
		"username": req.Username,
		"email":    req.Email,
		"role":     req.Role,
	}, "system.events.v1.staff.updated")

	resp := toResponse(*updated)
	return &resp, nil
}

func (s *service) Delete(ctx context.Context, ownerUserID, staffID uint) error {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, err)
		return err
	}

	existing, err := s.repo.FindByIDInCenter(ctx, centerID, staffID)
	if err != nil {
		s.publishError(ctx, types.ActionFailed, err)
		return err
	}
	if existing == nil {
		return ErrStaffNotFound
	}

	if err := s.repo.SoftDelete(ctx, staffID); err != nil {
		s.publishError(ctx, types.ActionFailed, err)
		return err
	}

	// Revoke any active sessions the removed staff member still holds.
	if err := s.tokenRevoker.DeleteTokensByUserID(staffID, "access"); err != nil {
		log.Printf("failed to revoke access tokens for deleted staff %d: %v", staffID, err)
	}
	if err := s.tokenRevoker.DeleteTokensByUserID(staffID, "refresh"); err != nil {
		log.Printf("failed to revoke refresh tokens for deleted staff %d: %v", staffID, err)
	}

	s.publishEvent(ctx, types.LevelInfo, types.ActionDeleted, fmt.Sprint(staffID), fmt.Sprint(ownerUserID), http.StatusOK, map[string]any{
		"username": existing.Username,
		"email":    existing.Email,
	}, "system.events.v1.staff.deleted")

	return nil
}

func toResponse(u auth.User) Response {
	centerID := uint(0)
	if u.CenterID != nil {
		centerID = *u.CenterID
	}
	return Response{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role.Name,
		CenterID:  centerID,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

// ==========================================
// INTERNAL HELPERS FOR PUBLISHING EVENTS
// (module-local, mirrors auth/service.go and notifications/service.go —
// this codebase does not have a shared publish-audit-event helper)
// ==========================================

func (s *service) publishError(ctx context.Context, action types.EventAction, err error) {
	if err == nil {
		return
	}
	payload := map[string]string{"error_message": err.Error()}
	s.publishEvent(ctx, types.LevelError, action, "unknown", "system", http.StatusInternalServerError, payload, "system.events.v1.staff.system_error")
}

func (s *service) publishEvent(
	ctx context.Context,
	level types.LogLevel,
	action types.EventAction,
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
		Entity:     "Staff",
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
