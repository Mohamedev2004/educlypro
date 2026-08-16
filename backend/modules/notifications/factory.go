package notifications

import (
	"encoding/json"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/datatypes"
)

var fakeTopics = []string{
	"user.registered",
	"user.welcome",
	"attendance.submitted",
	"grade.published",
	"session.reminder",
}

var fakeChannels = []string{
	"in_app",
	"email",
}

func NewFakeNotification(userID uint) Notification {
	isRead := gofakeit.Bool()
	createdAt := gofakeit.DateRange(time.Now().AddDate(0, 0, -30), time.Now())

	payloadBytes, _ := json.Marshal(map[string]any{
		"source": "seeder",
		"kind":   "demo",
	})

	notification := Notification{
		RequestID: gofakeit.UUID(),
		UserID:    userID,
		Topic:     gofakeit.RandomString(fakeTopics),
		Title:     gofakeit.Sentence(4),
		Body:      gofakeit.Sentence(12),
		Payload:   datatypes.JSON(payloadBytes),
		Channel:   gofakeit.RandomString(fakeChannels),
		IsRead:    isRead,
		CreatedAt: createdAt,
	}

	if isRead {
		readAt := createdAt.Add(time.Duration(gofakeit.Number(1, 360)) * time.Minute)
		notification.ReadAt = &readAt
	}

	return notification
}

func NewFakeNotifications(userIDs []uint, count int) []Notification {
	if len(userIDs) == 0 || count <= 0 {
		return nil
	}

	notifications := make([]Notification, 0, count)
	for i := 0; i < count; i++ {
		userID := userIDs[i%len(userIDs)]
		notifications = append(notifications, NewFakeNotification(userID))
	}

	return notifications
}
