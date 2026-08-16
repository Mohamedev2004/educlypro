package main

import (
	"educlypro/config"
	"educlypro/shared/database"
)

func main() {
	config.Load()
	database.Connect()
	database.Migrate()
	database.SeedAll(database.MainDB)
}
