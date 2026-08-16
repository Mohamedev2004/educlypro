package schedules

import (
	"log"
	"time"
)

// Run starts a background worker for a task
func Run(name string, interval time.Duration, task func() error) {
	go func() {
		// Run once on startup
		if err := task(); err != nil {
			log.Printf("Schedule [%s] initial run error: %v", name, err)
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := task(); err != nil {
				log.Printf("Schedule [%s] error: %v", name, err)
			}
		}
	}()
}
