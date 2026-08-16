package logs

import (
	"encoding/json"
	"log"

	"educlypro/shared/types" // Adjust to your actual module name

	"github.com/ThreeDotsLabs/watermill/message"
)

type Subscriber struct {
	repo Repository
}

func NewSubscriber(repo Repository) *Subscriber {
	return &Subscriber{repo: repo}
}

// ProcessLogEvent is the function triggered by Watermill
func (s *Subscriber) ProcessLogEvent(msg *message.Message) error {
	var event types.AuditEvent

	// 1. Unmarshal the envelope
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		log.Printf("Malformed audit event dropped: %v", err)
		return nil // Return nil so Watermill drops it instead of retrying forever
	}

	// 2. Convert the dynamic 'any' payload into raw JSON bytes for GORM
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		payloadBytes = []byte(`{}`) // Fallback to empty JSON object
	}

	// 3. Map to database model
	auditLog := &AuditLog{
		RequestID:  event.RequestID,
		Level:      string(event.Level),
		Entity:     event.Entity,
		EntityID:   event.EntityID,
		ActorID:    event.ActorID,
		StatusCode: event.StatusCode,
		Action:     string(event.Action),
		Payload:    payloadBytes,
		DurationMs: event.DurationMs,
		CreatedAt:  event.Timestamp,
	}

	// 4. Save using the background context from the message
	if err := s.repo.Create(msg.Context(), auditLog); err != nil {
		log.Printf("Failed to save audit log: %v", err)
		return err // Returning an error tells Watermill to retry this message later
	}

	return nil
}
