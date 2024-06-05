package database

import (
	"fmt"

	"github.com/davidalvarez305/budgeting/models"
	"github.com/davidalvarez305/budgeting/types"
)

func InsertCSRFToken(token models.CSRFToken) error {
	stmt, err := DB.Prepare("INSERT INTO csrf_token(expiry_time, token) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token.ExpiryTime, token.Token)
	if err != nil {
		return err
	}

	fmt.Println("CSRFToken inserted successfully")
	return nil
}

func GetCSRFToken(decryptedToken string) (models.CSRFToken, error) {
	var token models.CSRFToken

	stmt, err := DB.Prepare("SELECT * FROM csrf_token WHERE token = ?")
	if err != nil {
		return token, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(decryptedToken)

	err = row.Scan(&token.CSRFTokenID, &token.ExpiryTime, &token.Token, &token.IsUsed)
	if err != nil {
		return token, err
	}

	return token, nil
}

func CreateLeadAndMarketing(quoteForm types.QuoteForm) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Insert Lead
	leadQuery := `
		INSERT INTO lead (first_name, last_name, phone_number, created_at, rent, foot_traffic, foot_traffic_type, vending_type_id, vending_location_id, city_id)
		VALUES (?, ?, ?, UNIX_TIMESTAMP(), ?, ?, ?, ?, ?, ?)
	`
	leadResult, err := tx.Exec(leadQuery, quoteForm.FirstName, quoteForm.LastName, quoteForm.PhoneNumber, quoteForm.Rent, quoteForm.FootTraffic, quoteForm.FootTrafficType, quoteForm.MachineType, quoteForm.LocationType, quoteForm.City)
	if err != nil {
		return err
	}
	leadID, err := leadResult.LastInsertId()
	if err != nil {
		return err
	}

	// Insert Lead Marketing
	marketingQuery := `
		INSERT INTO lead_marketing (lead_id, source, medium, channel, landing_page, keyword, referrer, gclid, campaign_id, ad_campaign, ad_group_id, ad_group_name, ad_set_id, ad_set_name, ad_id, ad_headline, language, os, user_agent, button_clicked, device_type, ip)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(marketingQuery, leadID, quoteForm.Source, quoteForm.Medium, quoteForm.Channel, quoteForm.LandingPage, quoteForm.Keyword, quoteForm.Referrer, quoteForm.Gclid, quoteForm.CampaignID, quoteForm.AdCampaign, quoteForm.AdGroupID, quoteForm.AdGroupName, quoteForm.AdSetID, quoteForm.AdSetName, quoteForm.AdID, quoteForm.AdHeadline, quoteForm.Language, quoteForm.OS, quoteForm.UserAgent, quoteForm.ButtonClicked, quoteForm.DeviceType, quoteForm.IP)
	if err != nil {
		return err
	}

	return nil
}

func MarkCSRFTokenAsUsed(token models.CSRFToken) error {
	stmt, err := DB.Prepare("UPDATE csrf_token SET is_used = true WHERE token = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token.Token)
	if err != nil {
		return err
	}

	fmt.Println("CSRFToken marked as used successfully")
	return nil
}

func SaveSMS(msg models.TextMessage) error {
	query := `
		INSERT INTO text_message (message_sid, from_number, user_id, to_number, body, status, created_at, is_inbound)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := DB.Exec(query, msg.MessageSID, msg.UserID, msg.FromNumber, msg.ToNumber, msg.Body, msg.Status, msg.CreatedAt, msg.IsInbound)
	return err
}
