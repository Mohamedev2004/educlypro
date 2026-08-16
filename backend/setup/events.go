package setup

import (
	"context"
	"log"

	"educlypro/modules/logs"
	"educlypro/shared/database"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func EventBus() *gochannel.GoChannel {
	logger := watermill.NewStdLogger(false, false)
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	router := newRouter(logger, pubSub)

	go func() {
		if err := router.Run(context.Background()); err != nil {
			log.Fatalf("Event router stopped: %v", err)
		}
	}()
	<-router.Running()

	return pubSub
}

func newRouter(logger watermill.LoggerAdapter, pubSub *gochannel.GoChannel) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		log.Fatalf("Failed to create event router: %v", err)
	}

	logRepo := logs.NewRepository(database.LogsDB)
	router.AddNoPublisherHandler(
		"audit_log_writer",
		"system.audit_logs",
		pubSub,
		logs.NewSubscriber(logRepo).ProcessLogEvent,
	)

	RegisterNotifications(router, pubSub)

	return router
}
