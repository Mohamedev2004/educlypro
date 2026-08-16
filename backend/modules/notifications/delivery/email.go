package delivery

import (
	"educlypro/shared/types"
	"educlypro/shared/utils"
	"context"
	"log"
)

type emailDispatcher struct{}

func NewEmailDispatcher() Dispatcher {
	return &emailDispatcher{}
}

func (d *emailDispatcher) Channel() types.NotificationChannel {
	return types.ChannelEmail
}

func (d *emailDispatcher) Send(ctx context.Context, userID uint, email string, event *types.NotificationEvent) error {
	if email == "" {
		log.Printf("[EMAIL] Skipping send to UserID %d: email address missing", userID)
		return nil
	}

	log.Printf("[EMAIL] Sending real email to %s: %s", email, event.Title)

	switch event.Topic {
	case "user.welcome":
		username, _ := event.Payload["username"].(string)
		return utils.SendWelcomeEmail(email, username)
	case "user.registered":
		username, _ := event.Payload["username"].(string)
		userEmail, _ := event.Payload["email"].(string)
		return utils.SendAdminRegistrationEmail(email, username, userEmail)
	default:
		return utils.SendEmail(email, event.Title, event.Body)
	}
}
