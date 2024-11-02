package services

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/models"
)

func StartEmailScheduler() {
	go func() {
		for {
			now := time.Now()

			nextRun := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
			if now.Month() == 12 {
				nextRun = nextRun.AddDate(1, 0, 0) // Increment year if current month is December
			}
			timeUntilNextRun := nextRun.Sub(now)

			// Wait until the next scheduled run
			time.Sleep(timeUntilNextRun)

			scheduledEmails, err := database.GetScheduledEmails()
			if err != nil {
				log.Printf("Unable to retrieve scheduled emails: %v", err)
				continue
			}

			// Send emails based on the retrieved scheduled emails
			for _, email := range scheduledEmails {

				if now.Unix() > email.LastSent+email.IntervalSeconds {

					recipients := strings.Split(email.Recipients, ",")
					subject := email.Subject
					sender := email.Sender

					fileName := fmt.Sprintf("%s_%s_%s.xls", email.EmailName, now.Local().Month(), now.Local().Year())
					s3Key := "email-attachments/" + fileName
					localFilePath := constants.EMAIL_ATTACHMENTS_DIR + fileName

					// Find a way to save SQL queries so that they can be dynamically generated (( POSTGRES FUNCTIONS ))
					// Read file from AWS S3 SQL, save SQL file from
					data, err := database.GenerateEmailDataWithAttachment(email)
					if err != nil {
						continue
					}

					excelFilePath, err := helpers.GenerateExcelFile(data, "data", localFilePath)
					if err != nil {
						continue
					}

					fileInfo, err := os.Open(excelFilePath)
					if err != nil {
						continue
					}

					info, err := fileInfo.Stat()
					if err != nil {
						continue
					}

					err = UploadFileToS3(fileInfo, info.Size(), s3Key)
					if err != nil {
						continue
					}

					err = SendGmailWithAttachment(recipients, subject, sender, body, excelFilePath)
					if err != nil {
						continue
					}

					// Update LastSent to the current time after successful sending
					email.LastSent = now.Unix()

					err = database.UpdateEmailScheduledLastSent(email.EmailID, email.LastSent)
					if err != nil {
						log.Printf("Unable to update scheduled emails: %v", err)
						continue
					}

					sentEmail := models.SentEmail{
						EmailID:        email.EmailID,
						DeliveryStatus: "",
						DateSent:       now.Unix(),
						ErrorMessage:   "",
					}

					err = database.CreateScheduledEmailLog(sentEmail)
					if err != nil {
						log.Printf("Unable to create scheduled email log: %v", err)
						continue
					}
				}
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
