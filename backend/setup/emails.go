package setup

import (
	"educlypro/shared/types"
	"educlypro/shared/utils"
	"encoding/json"
	"errors"
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func RegisterEmails(router *message.Router, pubSub *gochannel.GoChannel) {
	router.AddNoPublisherHandler(
		"email_password_reset",
		"system.events.v1.auth.password_reset_requested",
		pubSub,
		handlePasswordResetEmail,
	)

	// Future emails — just add a new handler here:
	// router.AddNoPublisherHandler("email_welcome", "system.events.v1.auth.registered", pubSub, handleWelcomeEmail)
}

func handlePasswordResetEmail(msg *message.Message) error {
	payload, err := extractPayload(msg)
	if err != nil {
		return err
	}

	email, _ := payload["email"].(string)
	token, _ := payload["token"].(string)

	return utils.SendResetPasswordEmail(email, "", token)
}

// Future handlers:
// func handleWelcomeEmail(msg *message.Message) error { ... }

func extractPayload(msg *message.Message) (map[string]any, error) {
	var event types.AuditEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		log.Printf("Failed to unmarshal email event: %v", err)
		return nil, err
	}

	payload, ok := event.Payload.(map[string]any)
	if !ok {
		return nil, errors.New("invalid event payload")
	}

	return payload, nil
}
