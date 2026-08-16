package notifications

import (
	"educlypro/shared/types"
	"encoding/json"
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
)

type Subscriber struct {
	service Service
}

func NewSubscriber(service Service) *Subscriber {
	return &Subscriber{service: service}
}

func (s *Subscriber) ProcessNotificationEvent(msg *message.Message) error {
	var event types.NotificationEvent

	log.Printf("Received notification event: %s", msg.UUID)

	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		log.Printf("Malformed notification event dropped: %v", err)
		return nil // Drop malformed messages
	}

	if err := s.service.Dispatch(msg.Context(), &event); err != nil {
		log.Printf("Failed to dispatch notification: %v", err)
		return err // Retry on error
	}

	return nil
}
