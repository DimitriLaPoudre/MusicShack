package utils

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"
)

func Fetch(ctx context.Context) error {
	return nil
}

func AutoFetch(ctx context.Context) *cron.Cron {
	c := cron.New()

	c.AddFunc("* * * * *", func() {
		for try := range 3 {
			if err := Fetch(ctx); err != nil {
				log.Println("AutoFetch: try ", try, ": ", err)
			} else {
				log.Println("AutoFetch: success")
				break
			}
		}
	})

	c.Start()
	return c
}
