package notifications

import (
	"educlypro/modules/notifications/delivery"
	"educlypro/shared/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type UserResolver interface {
	UserIDsByRole(ctx context.Context, roles []string) ([]uint, error)
	UserEmailByID(ctx context.Context, userID uint) (string, error)
}

type Service interface {
	Dispatch(ctx context.Context, event *types.NotificationEvent) error
	List(ctx context.Context, userID uint, params ListParams) (*ListResponse, error)
	MarkRead(ctx context.Context, userID uint, id uint) error
	MarkAllRead(ctx context.Context, userID uint) error
	UnreadCount(ctx context.Context, userID uint) (int64, error)
	Delete(ctx context.Context, userID uint, id uint) error
	DeleteAllRead(ctx context.Context, userID uint) error
}

type service struct {
	repo        Repository
	resolver    UserResolver
	dispatchers map[types.NotificationChannel]delivery.Dispatcher
	publisher   message.Publisher
}

func NewService(repo Repository, resolver UserResolver, dispatchers []delivery.Dispatcher, publisher message.Publisher) Service {
	dMap := make(map[types.NotificationChannel]delivery.Dispatcher)
	for _, d := range dispatchers {
		dMap[d.Channel()] = d
	}
	return &service{
		repo:        repo,
		resolver:    resolver,
		dispatchers: dMap,
		publisher:   publisher,
	}
}

func (s *service) Dispatch(ctx context.Context, event *types.NotificationEvent) error {
	s.publishEvent(ctx, types.LevelInfo, types.ActionCreated, "Notification", event.Topic, "system", http.StatusOK, map[string]any{
		"title":    event.Title,
		"topic":    event.Topic,
		"channels": event.Channels,
	}, "system.events.v1.notifications.dispatch_started")

	// 1. Resolve recipients
	recipientSet := make(map[uint]struct{})
	for _, id := range event.RecipientIDs {
		recipientSet[id] = struct{}{}
	}

	if len(event.RoleTargets) > 0 {
		roleUserIDs, err := s.resolver.UserIDsByRole(ctx, event.RoleTargets)
		if err != nil {
			s.publishError(ctx, types.ActionFailed, "Notification", err)
			return fmt.Errorf("failed to resolve roles: %w", err)
		}
		for _, id := range roleUserIDs {
			recipientSet[id] = struct{}{}
		}
	}

	// 2. Fan-out per recipient and per channel
	successCount := 0
	failCount := 0

	for userID := range recipientSet {
		// Resolve user email for channels that need it (like email)
		userEmail, err := s.resolver.UserEmailByID(ctx, userID)
		if err != nil {
			s.publishEvent(ctx, types.LevelWarn, types.ActionFailed, "Notification", fmt.Sprint(userID), "system", http.StatusInternalServerError, map[string]any{
				"user_id": userID,
				"error":   err.Error(),
			}, "system.events.v1.notifications.resolver_error")
		}

		for _, channel := range event.Channels {
			// 3. Send per channel
			if d, ok := s.dispatchers[channel]; ok {
				if err := d.Send(ctx, userID, userEmail, event); err != nil {
					failCount++
					s.publishEvent(ctx, types.LevelError, types.ActionFailed, "Notification", fmt.Sprint(userID), "system", http.StatusInternalServerError, map[string]any{
						"user_id": userID,
						"channel": channel,
						"error":   err.Error(),
					}, "system.events.v1.notifications.delivery_failed")
					continue
				}
				successCount++
			}
		}
	}

	s.publishEvent(ctx, types.LevelInfo, types.ActionCreated, "Notification", event.Topic, "system", http.StatusOK, map[string]any{
		"total_recipients": len(recipientSet),
		"success_count":    successCount,
		"fail_count":       failCount,
	}, "system.events.v1.notifications.dispatch_completed")

	return nil
}

func (s *service) List(ctx context.Context, userID uint, params ListParams) (*ListResponse, error) {
	offset := (params.Page - 1) * params.PerPage

	items, err := s.repo.ListForUser(ctx, userID, params.Filter, params.PerPage, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountForUser(ctx, userID, params.Filter)
	if err != nil {
		return nil, err
	}

	counts, err := s.repo.CountSummary(ctx, userID)
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

	return &ListResponse{
		Items:  items,
		Filter: params.Filter,
		Counts: counts,
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

func (s *service) MarkRead(ctx context.Context, userID uint, id uint) error {
	if err := s.repo.MarkRead(ctx, userID, id); err != nil {
		s.publishError(ctx, types.ActionFailed, "Notification", err)
		return err
	}
	s.publishEvent(ctx, types.LevelInfo, types.ActionUpdated, "Notification", fmt.Sprint(id), fmt.Sprint(userID), http.StatusOK, nil, "system.events.v1.notifications.marked_read")
	return nil
}

func (s *service) MarkAllRead(ctx context.Context, userID uint) error {
	if err := s.repo.MarkAllRead(ctx, userID); err != nil {
		s.publishError(ctx, types.ActionFailed, "Notification", err)
		return err
	}
	s.publishEvent(ctx, types.LevelInfo, types.ActionUpdated, "Notification", "all", fmt.Sprint(userID), http.StatusOK, nil, "system.events.v1.notifications.all_marked_read")
	return nil
}

func (s *service) UnreadCount(ctx context.Context, userID uint) (int64, error) {
	return s.repo.UnreadCount(ctx, userID)
}

func (s *service) Delete(ctx context.Context, userID uint, id uint) error {
	if err := s.repo.Delete(ctx, userID, id); err != nil {
		s.publishError(ctx, types.ActionFailed, "Notification", err)
		return err
	}
	s.publishEvent(ctx, types.LevelInfo, types.ActionDeleted, "Notification", fmt.Sprint(id), fmt.Sprint(userID), http.StatusOK, nil, "system.events.v1.notifications.deleted")
	return nil
}

func (s *service) DeleteAllRead(ctx context.Context, userID uint) error {
	if err := s.repo.DeleteAllRead(ctx, userID); err != nil {
		s.publishError(ctx, types.ActionFailed, "Notification", err)
		return err
	}
	s.publishEvent(ctx, types.LevelInfo, types.ActionDeleted, "Notification", "all_read", fmt.Sprint(userID), http.StatusOK, nil, "system.events.v1.notifications.all_read_deleted")
	return nil
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
	s.publishEvent(ctx, types.LevelError, action, entity, "unknown", "system", http.StatusInternalServerError, payload, "system.events.v1.notifications.system_error")
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
