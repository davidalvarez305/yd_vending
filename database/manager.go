package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/models"
	"github.com/davidalvarez305/yd_vending/types"
)

func InsertCSRFToken(token models.CSRFToken) error {
	stmt, err := DB.Prepare(`INSERT INTO "csrf_token" ("expiry_time", "token", "is_used") VALUES(to_timestamp($1), $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token.ExpiryTime, token.Token, token.IsUsed)
	if err != nil {
		return err
	}

	fmt.Println("CSRFToken inserted successfully")
	return nil
}

func CheckIsTokenUsed(decryptedToken string) (bool, error) {
	var isUsed bool

	stmt, err := DB.Prepare(`SELECT is_used FROM "csrf_token" WHERE "token" = $1`)
	if err != nil {
		return isUsed, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(decryptedToken)

	err = row.Scan(&isUsed)
	if err != nil {
		return isUsed, err
	}

	return isUsed, nil
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

	var leadID int
	leadQuery := `
		INSERT INTO lead (first_name, last_name, phone_number, created_at, rent, foot_traffic, foot_traffic_type, vending_type_id, vending_location_id, city_id)
		VALUES ($1, $2, $3, to_timestamp($4), $5, $6, $7, $8, $9, $10)
		RETURNING lead_id
	`
	err = tx.QueryRow(leadQuery, quoteForm.FirstName, quoteForm.LastName, quoteForm.PhoneNumber, time.Now().Unix(), quoteForm.Rent, quoteForm.FootTraffic, quoteForm.FootTrafficType, quoteForm.MachineType, quoteForm.LocationType, quoteForm.City).Scan(&leadID)
	if err != nil {
		fmt.Println("ERROR INSERTING LEAD")
		return err
	}

	marketingQuery := `
		INSERT INTO lead_marketing (lead_id, source, medium, channel, landing_page, keyword, referrer, click_id, campaign_id, ad_campaign, ad_group_id, ad_group_name, ad_set_id, ad_set_name, ad_id, ad_headline, language, user_agent, button_clicked, ip, google_user_id, google_client_id, csrf_secret, facebook_click_id, facebook_client_id, longitude, latitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)
	`
	_, err = tx.Exec(marketingQuery, leadID, quoteForm.Source, quoteForm.Medium, quoteForm.Channel, quoteForm.LandingPage, quoteForm.Keyword, quoteForm.Referrer, quoteForm.ClickID, quoteForm.CampaignID, quoteForm.AdCampaign, quoteForm.AdGroupID, quoteForm.AdGroupName, quoteForm.AdSetID, quoteForm.AdSetName, quoteForm.AdID, quoteForm.AdHeadline, quoteForm.Language, quoteForm.UserAgent, quoteForm.ButtonClicked, quoteForm.IP, quoteForm.GoogleUserID, quoteForm.GoogleClientID, quoteForm.CSRFSecret, quoteForm.FacebookClickID, quoteForm.FacebookClientID, quoteForm.Longitude, quoteForm.Latitude)
	if err != nil {
		fmt.Println("ERROR INSERTING MARKETING")
		return err
	}

	return nil
}

func MarkCSRFTokenAsUsed(token string) error {
	stmt, err := DB.Prepare(`UPDATE "csrf_token" SET "is_used" = true WHERE "token" = $1`)
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

func SaveSMS(msg models.Message) error {
	query := `
		INSERT INTO message (external_id, user_id, lead_id, text, date_created, text_from, text_to, is_inbound)
		VALUES ($1, $2, $3, $4, to_timestamp($5), $6, $7, $8)
	`
	_, err := DB.Exec(query, msg.ExternalID, msg.UserID, msg.LeadID, msg.Text, msg.DateCreated, msg.TextFrom, msg.TextTo, msg.IsInbound)
	return err
}

func SavePhoneCall(phoneCall models.PhoneCall) error {
	query := `
		INSERT INTO phone_call (
			external_id, user_id, lead_id, call_duration,
			date_created, call_from, call_to, is_inbound,
			recording_url, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	err := DB.QueryRow(
		query,
		phoneCall.ExternalID,
		phoneCall.UserID,
		phoneCall.LeadID,
		phoneCall.CallDuration,
		phoneCall.DateCreated,
		phoneCall.CallFrom,
		phoneCall.CallTo,
		phoneCall.IsInbound,
		phoneCall.RecordingURL,
		phoneCall.Status,
	)

	if err != nil {
		return fmt.Errorf("error saving phone call: %w", err)
	}

	return nil
}

func GetUserIDFromPhoneNumber(from string) (int, error) {
	var userId int

	stmt, err := DB.Prepare(`SELECT "user_id" FROM "user" WHERE "phone_number" = $1`)
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

func GetPhoneNumberFromUserID(userID int) (string, error) {
	var phoneNumber string

	stmt, err := DB.Prepare(`SELECT "phone_number" FROM "user" WHERE "user_id" = $1`)
	if err != nil {
		return phoneNumber, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID)

	err = row.Scan(&phoneNumber)
	if err != nil {
		return phoneNumber, err
	}

	return phoneNumber, nil
}

func GetUserById(id int) (models.User, error) {
	var user models.User

	stmt, err := DB.Prepare(`SELECT * FROM "user" WHERE "user_id" = $1`)
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	err = row.Scan(&user.UserID, &user.Email, &user.Password, &user.IsAdmin, &user.PhoneNumber, &user.FirstName, &user.LastName)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	stmt, err := DB.Prepare(`SELECT * FROM "user" WHERE "email" = $1`)
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(email)

	err = row.Scan(&user.UserID, &user.Email, &user.Password, &user.IsAdmin, &user.PhoneNumber, &user.FirstName, &user.LastName)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetVendingTypes() ([]models.VendingType, error) {
	var vendingTypes []models.VendingType

	rows, err := DB.Query(`SELECT * FROM "vending_type"`)
	if err != nil {
		return vendingTypes, err
	}
	defer rows.Close()

	for rows.Next() {
		var vt models.VendingType
		err := rows.Scan(&vt.VendingTypeID, &vt.MachineType)
		if err != nil {
			return vendingTypes, err
		}
		vendingTypes = append(vendingTypes, vt)
	}

	if err := rows.Err(); err != nil {
		return vendingTypes, err
	}

	return vendingTypes, nil
}

func GetVendingLocations() ([]models.VendingLocation, error) {
	var vendingLocations []models.VendingLocation

	rows, err := DB.Query(`SELECT * FROM "vending_location"`)
	if err != nil {
		return vendingLocations, err
	}
	defer rows.Close()

	for rows.Next() {
		var vl models.VendingLocation
		err := rows.Scan(&vl.VendingLocationID, &vl.LocationType)
		if err != nil {
			return vendingLocations, err
		}
		vendingLocations = append(vendingLocations, vl)
	}

	if err := rows.Err(); err != nil {
		return vendingLocations, err
	}

	return vendingLocations, nil
}

func GetCities() ([]models.City, error) {
	var cities []models.City

	rows, err := DB.Query(`SELECT "city_id", "name" FROM "city"`)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return cities, err
	}
	defer rows.Close()

	for rows.Next() {
		var city models.City
		err := rows.Scan(&city.CityID, &city.Name)
		if err != nil {
			return cities, err
		}
		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
		return cities, err
	}

	return cities, nil
}

func GetLeadList(params types.GetLeadsParams) ([]types.LeadList, int, error) {
	var leads []types.LeadList

	query := `SELECT l.lead_id, l.first_name, l.last_name, l.phone_number, 
		l.created_at, l.rent, l.foot_traffic, l.foot_traffic_type, 
		vt.machine_type, vl.location_type, c.name as city, lm.language,
		l.city_id, l.vending_type_id, l.vending_location_id,
		COUNT(*) OVER() AS total_rows
		FROM lead AS l
		JOIN city AS c ON c.city_id = l.city_id
		JOIN vending_type AS vt ON vt.vending_type_id = l.vending_type_id
		JOIN vending_location AS vl ON vl.vending_location_id = l.vending_location_id
		JOIN lead_marketing AS lm ON lm.lead_id = l.lead_id
		WHERE 1=1`

	args := []interface{}{}
	argIdx := 1

	// Add conditions based on non-empty fields in params
	if params.VendingType != "" {
		query += fmt.Sprintf(" AND vt.vending_type_id = $%d", argIdx)
		args = append(args, params.VendingType)
		argIdx++
	}

	if params.LocationType != "" {
		query += fmt.Sprintf(" AND vl.vending_location_id = $%d", argIdx)
		args = append(args, params.LocationType)
		argIdx++
	}

	if params.City != "" {
		query += fmt.Sprintf(" AND c.city_id = $%d", argIdx)
		args = append(args, params.City)
		argIdx++
	}

	query += fmt.Sprintf(" LIMIT $%d", argIdx)
	args = append(args, constants.LeadsPerPage)
	argIdx++

	if params.PageNum != "" {
		pageNum, err := strconv.Atoi(params.PageNum)
		if err != nil {
			fmt.Printf("Could not convert page num: %+v\n", err)
			return leads, 0, err
		}

		offset := (pageNum - 1) * constants.LeadsPerPage
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, offset)
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return leads, 0, err
	}
	defer rows.Close()

	var totalRows int
	for rows.Next() {
		var lead types.LeadList
		var createdAt time.Time

		err := rows.Scan(&lead.LeadID,
			&lead.FirstName,
			&lead.LastName,
			&lead.PhoneNumber,
			&createdAt,
			&lead.Rent,
			&lead.FootTraffic,
			&lead.FootTrafficType,
			&lead.MachineType,
			&lead.LocationType,
			&lead.City,
			&lead.Language,
			&lead.CityID,
			&lead.VendingTypeID,
			&lead.VendingLocationID,
			&totalRows)
		if err != nil {
			return leads, 0, err
		}
		lead.CreatedAt = createdAt.Unix()
		leads = append(leads, lead)
	}

	if err := rows.Err(); err != nil {
		return leads, 0, err
	}

	return leads, totalRows, nil
}

func GetLeadDetails(leadID string) (types.LeadDetails, error) {
	query := `SELECT l.lead_id,
	l.first_name,
	l.last_name,
	l.phone_number,
	vt.machine_type,
	vl.location_type,
	lm.ad_campaign,
	lm.medium,
	lm.source,
	lm.referrer,
	lm.landing_page,
	lm.ip,
	lm.keyword,
	lm.channel,
	lm.language,
	c."name"
	FROM lead l
	JOIN vending_type vt ON l.vending_type_id = vt.vending_type_id
	JOIN vending_location vl ON l.vending_location_id = vl.vending_location_id
	JOIN lead_marketing lm ON l.lead_id = lm.lead_id
	JOIN city c ON c.city_id = l.city_id
	WHERE l.lead_id = $1`

	var leadDetails types.LeadDetails

	// Execute the query
	row := DB.QueryRow(query, leadID)

	// Scan the result into the LeadDetails struct
	err := row.Scan(
		&leadDetails.LeadID,
		&leadDetails.FirstName,
		&leadDetails.LastName,
		&leadDetails.PhoneNumber,
		&leadDetails.VendingType,
		&leadDetails.VendingLocation,
		&leadDetails.CampaignName,
		&leadDetails.Medium,
		&leadDetails.Source,
		&leadDetails.Referrer,
		&leadDetails.LandingPage,
		&leadDetails.IP,
		&leadDetails.Keyword,
		&leadDetails.Channel,
		&leadDetails.Language,
		&leadDetails.City,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return leadDetails, fmt.Errorf("no lead found with ID %d", leadID)
		}
		return leadDetails, err
	}

	return leadDetails, nil
}

func GetLeadIDFromPhoneNumber(from string) (int, error) {
	var leadId int

	stmt, err := DB.Prepare(`SELECT "lead_id" FROM "lead" WHERE "phone_number" = $1`)
	if err != nil {
		return leadId, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(from)

	err = row.Scan(&leadId)
	if err != nil {
		return leadId, err
	}

	return leadId, nil
}

func GetMessagesByLeadID(leadId int) ([]types.FrontendMessage, error) {
	var messages []types.FrontendMessage

	query := `SELECT CONCAT(l.first_name, ' ', l.last_name) as client_name,
	CONCAT(u.first_name, ' ', u.last_name) as user_name,
	m.text,
	m.date_created,
	m.is_inbound
	FROM "message" AS m
	JOIN "lead" AS l ON l.lead_id  = m.lead_id 
	JOIN "user" AS u ON u.user_id = m.user_id
	WHERE m.lead_id = $1
	ORDER BY m.date_created DESC`

	rows, err := DB.Query(query, leadId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return messages, err
	}
	defer rows.Close()

	for rows.Next() {
		var dateCreated time.Time

		var message types.FrontendMessage
		err := rows.Scan(
			&message.ClientName,
			&message.UserName,
			&message.Message,
			&dateCreated,
			&message.IsInbound,
		)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return messages, err
		}

		message.DateCreated = dateCreated.Unix()
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return messages, err
	}

	return messages, nil
}

func UpdateLead(form types.UpdateLeadForm) error {
	query := `
		UPDATE lead
		SET first_name = $2, last_name = $3, phone_number = $4, city_id = $5, vending_type_id = $6, vending_location_id = $7
		WHERE lead_id = $1
	`

	_, err := DB.Exec(query, form.LeadID, form.FirstName, form.LastName, form.PhoneNumber, form.City, form.VendingType, form.VendingLocation)
	if err != nil {
		return fmt.Errorf("failed to update lead: %v", err)
	}

	return nil
}

func UpdateLeadMarketing(form types.UpdateLeadMarketingForm) error {
	query := `
		UPDATE lead_marketing
		SET campaign_name = $2, medium = $3, source = $4, referrer = $5, landing_page = $6,
			ip = $7, keyword = $8, channel = $9, language = $10
		WHERE lead_id = $1
	`

	_, err := DB.Exec(query, form.LeadID, form.CampaignName, form.Medium, form.Source, form.Referrer,
		form.LandingPage, form.IP, form.Keyword, form.Channel, form.Language)
	if err != nil {
		return fmt.Errorf("failed to update lead marketing: %v", err)
	}

	return nil
}

func GetForwardPhoneNumber(to, from string) (types.IncomingPhoneCallForwarding, error) {
	var forwardingCall types.IncomingPhoneCallForwarding

	stmt, err := DB.Prepare(`SELECT u.first_name, u.user_id FROM "user" AS u WHERE u.phone_number = $1`)
	if err != nil {
		return forwardingCall, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(to)

	err = row.Scan(&forwardingCall.FirstName, &forwardingCall.UserID)
	if err != nil {
		return forwardingCall, err
	}

	stmt, err = DB.Prepare(`SELECT l.lead_id FROM "lead" AS l WHERE l.phone_number = $1`)
	if err != nil {
		return forwardingCall, err
	}
	defer stmt.Close()

	row = stmt.QueryRow(from)

	var leadID sql.NullInt64
	err = row.Scan(&leadID)
	if err != nil && err != sql.ErrNoRows {
		return forwardingCall, err
	}

	if leadID.Valid {
		forwardingCall.LeadID = int(leadID.Int64)
	} else {
		forwardingCall.LeadID = 0
	}

	switch forwardingCall.FirstName {
	case "Yovana":
		forwardingCall.ForwardPhoneNumber = "+1" + constants.YovaPhoneNumber
	case "David":
		forwardingCall.ForwardPhoneNumber = "+1" + constants.DavidPhoneNumber
	default:
		return forwardingCall, errors.New("no matching phone number")
	}

	return forwardingCall, nil
}

func GetPhoneCallBySID(sid string) (models.PhoneCall, error) {
	var phoneCall models.PhoneCall

	stmt, err := DB.Prepare(`SELECT * FROM phone_call WHERE external_id = $1`)
	if err != nil {
		return phoneCall, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(sid)

	err = row.Scan(
		&phoneCall.PhoneCallID,
		&phoneCall.ExternalID,
		&phoneCall.UserID,
		&phoneCall.LeadID,
		&phoneCall.CallDuration,
		&phoneCall.DateCreated,
		&phoneCall.CallFrom,
		&phoneCall.CallTo,
		&phoneCall.IsInbound,
		&phoneCall.RecordingURL,
		&phoneCall.Status,
	)
	if err != nil {
		return phoneCall, err
	}

	return phoneCall, nil
}

func UpdatePhoneCall(phoneCall models.PhoneCall) error {
	query := `
		UPDATE phone_call SET
			user_id = $1,
			lead_id = $2,
			call_duration = $3,
			date_created = $4,
			call_from = $5,
			call_to = $6,
			is_inbound = $7,
			recording_url = $8,
			status = $9
		WHERE external_id = $10`

	_, err := DB.Exec(
		query,
		phoneCall.UserID,
		phoneCall.LeadID,
		phoneCall.CallDuration,
		phoneCall.DateCreated,
		phoneCall.CallFrom,
		phoneCall.CallTo,
		phoneCall.IsInbound,
		phoneCall.RecordingURL,
		phoneCall.Status,
		phoneCall.ExternalID,
	)

	if err != nil {
		return fmt.Errorf("error updating phone call: %w", err)
	}

	return nil
}

func GetSession(userKey string) (models.Session, error) {
	var session models.Session
	sqlStatement := `
        SELECT session_id, user_id, csrf_secret, google_user_id, google_client_id, facebook_click_id, facebook_client_id, date_created, date_expires
        FROM sessions
        WHERE csrf_secret = $1
    `
	row := DB.QueryRow(sqlStatement, userKey)

	var dateCreated time.Time
	var dateExpires time.Time

	var userID sql.NullInt32
	err := row.Scan(&session.SessionID, &userID, &session.CSRFSecret, &session.GoogleUserID, &session.GoogleClientID, &session.FacebookClickID, &session.FacebookClientID, &dateCreated, &dateExpires)
	if err != nil {
		return session, err
	}

	if userID.Valid {
		session.UserID = int(userID.Int32)
	}

	session.DateCreated = dateCreated.Unix()
	session.DateExpires = dateExpires.Unix()

	return session, nil
}

func CreateSession(session models.Session) error {
	sqlStatement := `
        INSERT INTO sessions (csrf_secret, google_user_id, google_client_id, facebook_click_id, facebook_client_id, date_created, date_expires)
        VALUES ($1, $2, $3, $4, $5, to_timestamp($6), to_timestamp($7))
    `
	err := DB.QueryRow(sqlStatement, session.CSRFSecret, session.GoogleUserID, session.GoogleClientID, session.FacebookClickID, session.FacebookClientID, session.DateCreated, session.DateExpires)
	if err != nil {
		return err.Err()
	}

	return nil
}

func UpdateSession(session models.Session) error {
	sqlStatement := `
        UPDATE sessions
        SET google_user_id = $1,
            google_client_id = $2,
            facebook_click_id = $3,
            facebook_client_id = $4,
			user_id = $6
        WHERE csrf_secret = $5
    `
	_, err := DB.Exec(sqlStatement, session.GoogleUserID, session.GoogleClientID, session.FacebookClickID, session.FacebookClientID, session.CSRFSecret, session.UserID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSession(secret string) error {
	sqlStatement := `
        DELETE FROM sessions WHERE csrf_secret = $1
    `
	_, err := DB.Exec(sqlStatement, secret)
	if err != nil {
		return err
	}

	return nil
}
