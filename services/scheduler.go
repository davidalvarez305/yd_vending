package services

import (
	"log"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
)

func StartEmailScheduler() {
	go func() {
		for {
			now := time.Now()
			nextRun := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
			if now.Month() == 12 {
				nextRun = nextRun.AddDate(1, 0, 0)
			}
			timeUntilNextRun := nextRun.Sub(now)

			time.Sleep(timeUntilNextRun)

			if err := sendMonthlyEmail(); err != nil {
				log.Printf("Failed to send monthly email: %+v\n", err)

				// SAVE IN EMAIL SCHEDULE LOG DB
			}
		}
	}()
}

func sendMonthlyEmail() error {
	recipients := []string{}
	subject := "Monthly Commission Report"
	sender := constants.CompanyEmail
	body := "LINK TO MONTHLY REPORT"

	return SendGmail(recipients, subject, sender, body)
}
