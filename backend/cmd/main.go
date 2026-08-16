package main

import (
	"educlypro/config"
	"educlypro/routes"
	"educlypro/schedules"
	"educlypro/setup"
	"educlypro/shared/database"
	"context"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	database.Connect()
	database.Migrate()

	pubSub := setup.EventBus()
	schedules.StartAll(database.MainDB)

	// Build the Watermill router
	watermillRouter, err := message.NewRouter(message.RouterConfig{}, watermill.NewStdLogger(false, false))
	if err != nil {
		log.Fatal(err)
	}

	// Register ALL event handlers here
	setup.RegisterEmails(watermillRouter, pubSub)
	setup.RegisterNotifications(watermillRouter, pubSub)

	// Run the router in the background
	go func() {
		if err := watermillRouter.Run(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	r := gin.Default()
	routes.Register(r, database.MainDB, pubSub)
	r.Run(":" + config.Cfg.Server.Port)
}
