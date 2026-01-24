package cron

import (
	"log"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
)

func StartCronJobs() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(5).Minutes().Do(func() {
		log.Println("Pinging the health check endpoint...")
		resp, err := http.Get("http://localhost:8080/health")
		if err != nil {
			log.Printf("Error pinging health check endpoint: %v", err)
			return
		}
		defer resp.Body.Close()
		log.Println("Health check ping successful")
	})

	s.StartAsync()
}
