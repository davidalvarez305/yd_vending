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
	"github.com/davidalvarez305/yd_vending/types"
)

const (
	SuccessDeliveryStatus string = "SUCCESS"
	FailedDeliveryStatus  string = "FAILED"
)

var EmailTemplateFilePath string = constants.PARTIAL_TEMPLATES_DIR + "email_schedule_wrapper"

func StartEmailScheduler() {
	go func() {
		for {
			now := time.Now()
			loc, err := time.LoadLocation(constants.TimeZone)
			if err != nil {
				fmt.Println("Error loading location:", err)
				continue
			}

			nextRun := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, loc)
			if now.Month() == 12 {
				nextRun = nextRun.AddDate(1, 0, 0) // Increment year if current month is December
			}
			timeUntilNextRun := nextRun.Sub(now)

			time.Sleep(timeUntilNextRun)

			scheduledEmails, err := database.GetScheduledEmails()
			if err != nil {
				log.Printf("Unable to retrieve scheduled emails: %v", err)
				continue
			}

			for _, email := range scheduledEmails {
				if now.Unix() > email.LastSent+email.IntervalSeconds {
					recipients := strings.Split(email.Recipients, ", ")
					subject := email.Subject
					sender := email.Sender

					fileName := fmt.Sprintf("%s_%s_%d.xls", email.EmailName, now.Local().Month().String(), now.Local().Year())
					uploadReportS3Key := constants.EMAIL_ATTACHMENTS_S3_BUCKET + fileName
					localFilePath := constants.LOCAL_FILES_DIR + fileName

					sqlFileName := fmt.Sprintf("%s.sql", email.EmailName)
					sqlFileS3Key := constants.SQL_FILES_S3_BUCKET + sqlFileName
					sqlFileLocalPath := constants.LOCAL_FILES_DIR + sqlFileName

					sqlFile, err := DownloadFileFromS3(sqlFileS3Key, sqlFileLocalPath)
					if err != nil {
						continue
					}

					sqlQuery, err := os.ReadFile(sqlFile)
					if err != nil {
						continue
					}

					data, err := database.ExecuteQueryFromSQLFile(string(sqlQuery))
					if err != nil {
						continue
					}

					template, err := helpers.InsertHTMLIntoEmailTemplate(EmailTemplateFilePath, "content.html", email.Body, data)
					if err != nil {
						fmt.Printf("ERROR BUILDING SCHEDULED EMAIL TEMPLATE: %+v\n", err)
						continue
					}
					body := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n%s", template)

					excelFilePath, err := helpers.GenerateExcelFile(data, "data", localFilePath)
					if err != nil {
						log.Printf("Error generating Excel file: %v", err)
						continue
					}

					fileInfo, err := os.Open(excelFilePath)
					if err != nil {
						log.Printf("Error opening Excel file: %v", err)
						continue
					}

					func() {
						defer fileInfo.Close()

						info, err := fileInfo.Stat()
						if err != nil {
							log.Printf("Error getting file info: %v", err)
							return
						}

						err = UploadFileToS3(fileInfo, info.Size(), uploadReportS3Key)
						if err != nil {
							log.Printf("Error uploading file to S3: %v", err)
							return
						}

						timeSent := now.Unix()
						err = SendGmailWithAttachment(recipients, subject, sender, body, excelFilePath)
						if err != nil {
							// LOG FAILED DELIVERY
							deliveryStatus := FailedDeliveryStatus
							errorMessage := err.Error()
							sentEmail := types.SentEmail{
								EmailScheduleID: &email.EmailScheduleID,
								DeliveryStatus:  &deliveryStatus,
								DateSent:        &timeSent,
								ErrorMessage:    &errorMessage,
							}

							err = database.CreateSentEmail(sentEmail)
							if err != nil {
								log.Printf("Unable to create scheduled email log: %v", err)
							}
						} else {
							updateEmail := types.EmailScheduleForm{
								LastSent: &timeSent,
							}

							err = database.UpdateEmailSchedule(email.EmailScheduleID, updateEmail)
							if err != nil {
								log.Printf("Unable to update scheduled emails: %v", err)
							}

							// LOG SUCCESSFUL DELIVERY
							deliveryStatus := SuccessDeliveryStatus
							sentEmail := types.SentEmail{
								EmailScheduleID: &email.EmailScheduleID,
								DeliveryStatus:  &deliveryStatus,
								DateSent:        &timeSent,
								ErrorMessage:    nil,
							}

							err = database.CreateSentEmail(sentEmail)
							if err != nil {
								log.Printf("Unable to create scheduled email log: %v", err)
							}
						}
					}()
				}
			}
		}
	}()
}
