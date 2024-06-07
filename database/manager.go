package database

import (
	"fmt"

	"github.com/davidalvarez305/yd_vending/models"
	"github.com/davidalvarez305/yd_vending/types"
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

func CreateLeadAndMarketing(quoteForm types.QuoteForm, userKey []byte) error {
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
		INSERT INTO lead (first_name, last_name, phone_number, created_at, rent, foot_traffic, foot_traffic_type, vending_type_id, vending_location_id, city_id, user_key)
		VALUES (?, ?, ?, UNIX_TIMESTAMP(), ?, ?, ?, ?, ?, ?, ?)
	`
	leadResult, err := tx.Exec(leadQuery, quoteForm.FirstName, quoteForm.LastName, quoteForm.PhoneNumber, quoteForm.Rent, quoteForm.FootTraffic, quoteForm.FootTrafficType, quoteForm.MachineType, quoteForm.LocationType, quoteForm.City, string(userKey))
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

func MarkCSRFTokenAsUsed(token string) error {
	stmt, err := DB.Prepare("UPDATE csrf_token SET is_used = true WHERE token = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token)
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

func GetUserIDFromPhoneNumber(from string) (int, error) {
	var userId int

	stmt, err := DB.Prepare("SELECT user_id FROM user WHERE phone_number = ?")
	if err != nil {
		return userId, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(from)

	err = row.Scan(&userId)
	if err != nil {
		return userId, err
	}

	return userId, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	stmt, err := DB.Prepare("SELECT * FROM user WHERE email = ?")
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(email)

	err = row.Scan(&user.UserID, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetVendingTypes() ([]models.VendingType, error) {
	var vendingTypes []models.VendingType

	rows, err := DB.Query("SELECT vending_type_id, machine_type FROM vending_type;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vt models.VendingType
		err := rows.Scan(&vt.VendingTypeID, &vt.MachineType)
		if err != nil {
			return nil, err
		}
		vendingTypes = append(vendingTypes, vt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vendingTypes, nil
}

func GetVendingLocations() ([]models.VendingLocation, error) {
	var vendingLocations []models.VendingLocation

	rows, err := DB.Query("SELECT * FROM vending_location;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vl models.VendingLocation
		err := rows.Scan(&vl.VendingLocationID, &vl.LocationType)
		if err != nil {
			return nil, err
		}
		vendingLocations = append(vendingLocations, vl)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vendingLocations, nil
}

func GetCities() ([]models.City, error) {
	var cities []models.City

	rows, err := DB.Query("SELECT * FROM city;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var city models.City
		err := rows.Scan(&city.CityID, &city.Name)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cities, nil
}
