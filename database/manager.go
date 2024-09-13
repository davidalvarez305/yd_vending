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
	"github.com/davidalvarez305/yd_vending/utils"
)

func InsertCSRFToken(token models.CSRFToken) error {
	stmt, err := DB.Prepare(`INSERT INTO "csrf_token" ("expiry_time", "token", "is_used") VALUES(to_timestamp($1), $2, $3)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(token.ExpiryTime, token.Token, token.IsUsed)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CheckIsTokenUsed(decryptedToken string) (bool, error) {
	var isUsed bool

	stmt, err := DB.Prepare(`SELECT is_used FROM "csrf_token" WHERE "token" = $1`)
	if err != nil {
		return isUsed, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(decryptedToken)

	err = row.Scan(&isUsed)
	if err != nil {
		return isUsed, fmt.Errorf("error scanning row: %w", err)
	}

	return isUsed, nil
}

func CreateLeadAndMarketing(quoteForm types.QuoteForm) (int, error) {
	var leadID int
	tx, err := DB.Begin()
	if err != nil {
		return leadID, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()
	leadStmt, err := tx.Prepare(`
		INSERT INTO lead (first_name, last_name, phone_number, created_at, rent, foot_traffic, foot_traffic_type, vending_type_id, vending_location_id, message)
		VALUES ($1, $2, $3, to_timestamp($4), $5, $6, $7, $8, $9, $10)
		RETURNING lead_id
	`)
	if err != nil {
		return leadID, fmt.Errorf("error preparing lead statement: %w", err)
	}
	defer leadStmt.Close()

	err = leadStmt.QueryRow(quoteForm.FirstName, quoteForm.LastName, quoteForm.PhoneNumber, time.Now().Unix(), quoteForm.Rent, quoteForm.FootTraffic, quoteForm.FootTrafficType, quoteForm.MachineType, quoteForm.LocationType, quoteForm.Message).Scan(&leadID)
	if err != nil {
		return leadID, fmt.Errorf("error inserting lead: %w", err)
	}

	marketingStmt, err := tx.Prepare(`
		INSERT INTO lead_marketing (lead_id, source, medium, channel, landing_page, keyword, referrer, click_id, campaign_id, ad_campaign, ad_group_id, ad_group_name, ad_set_id, ad_set_name, ad_id, ad_headline, language, user_agent, button_clicked, ip, external_id, google_client_id, csrf_secret, facebook_click_id, facebook_client_id, longitude, latitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)
	`)
	if err != nil {
		return leadID, fmt.Errorf("error preparing marketing statement: %w", err)
	}
	defer marketingStmt.Close()

	_, err = marketingStmt.Exec(leadID, quoteForm.Source, quoteForm.Medium, quoteForm.Channel, quoteForm.LandingPage, quoteForm.Keyword, quoteForm.Referrer, quoteForm.ClickID, quoteForm.CampaignID, quoteForm.AdCampaign, quoteForm.AdGroupID, quoteForm.AdGroupName, quoteForm.AdSetID, quoteForm.AdSetName, quoteForm.AdID, quoteForm.AdHeadline, quoteForm.Language, quoteForm.UserAgent, quoteForm.ButtonClicked, quoteForm.IP, quoteForm.ExternalID, quoteForm.GoogleClientID, quoteForm.CSRFSecret, quoteForm.FacebookClickID, quoteForm.FacebookClientID, quoteForm.Longitude, quoteForm.Latitude)
	if err != nil {
		return leadID, fmt.Errorf("error inserting marketing data: %w", err)
	}

	err = tx.Commit()

	if err != nil {
		return leadID, fmt.Errorf("error committing transaction: %w", err)
	}

	return leadID, nil
}

func MarkCSRFTokenAsUsed(token string) error {
	stmt, err := DB.Prepare(`UPDATE "csrf_token" SET "is_used" = true WHERE "token" = $1`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(token)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	fmt.Println("CSRFToken marked as used successfully")
	return nil
}

func SaveSMS(msg models.Message) error {
	stmt, err := DB.Prepare(`
		INSERT INTO message (external_id, user_id, lead_id, text, date_created, text_from, text_to, is_inbound)
		VALUES ($1, $2, $3, $4, to_timestamp($5), $6, $7, $8)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var leadID sql.NullInt64
	if msg.LeadID != 0 {
		leadID = sql.NullInt64{Int64: int64(msg.LeadID), Valid: true}
	}

	_, err = stmt.Exec(msg.ExternalID, msg.UserID, leadID, msg.Text, msg.DateCreated, msg.TextFrom, msg.TextTo, msg.IsInbound)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func SavePhoneCall(phoneCall models.PhoneCall) error {
	stmt, err := DB.Prepare(`
		INSERT INTO phone_call (
			external_id, user_id, lead_id, call_duration,
			date_created, call_from, call_to, is_inbound,
			recording_url, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var leadID, callDuration sql.NullInt64

	if phoneCall.LeadID != 0 {
		leadID = sql.NullInt64{Int64: int64(phoneCall.LeadID), Valid: true}
	}
	if phoneCall.CallDuration != 0 {
		callDuration = sql.NullInt64{Int64: int64(phoneCall.CallDuration), Valid: true}
	}

	_, err = stmt.Exec(
		phoneCall.ExternalID,
		phoneCall.UserID,
		leadID,
		callDuration,
		phoneCall.DateCreated,
		phoneCall.CallFrom,
		phoneCall.CallTo,
		phoneCall.IsInbound,
		phoneCall.RecordingURL,
		phoneCall.Status,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetUserIDFromPhoneNumber(from string) (int, error) {
	var userId int

	stmt, err := DB.Prepare(`SELECT "user_id" FROM "user" WHERE "phone_number" = $1`)
	if err != nil {
		return userId, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(from)

	err = row.Scan(&userId)
	if err != nil {
		return userId, fmt.Errorf("error scanning row: %w", err)
	}

	return userId, nil
}

func GetConversionLeadInfo(leadId int) (types.ConversionLeadInfo, error) {
	var leadConversionInfo types.ConversionLeadInfo

	stmt, err := DB.Prepare(`SELECT l.lead_id, l.created_at, vt.machine_type, vl.location_type
		FROM "lead" AS l
	JOIN vending_type  AS vt ON vt.vending_type_id = l.vending_type_id
	JOIN vending_location AS vl ON vl.vending_location_id  = l.vending_location_id 
	WHERE l.lead_id = $1;`)

	if err != nil {
		return leadConversionInfo, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(leadId)

	var createdAt time.Time
	err = row.Scan(&leadConversionInfo.LeadID,
		&createdAt,
		&leadConversionInfo.MachineType,
		&leadConversionInfo.LocationType,
	)
	if err != nil {
		return leadConversionInfo, fmt.Errorf("error scanning row: %w", err)
	}

	leadConversionInfo.CreatedAt = createdAt.Unix()

	return leadConversionInfo, nil
}

func GetPhoneNumberFromUserID(userID int) (string, error) {
	var phoneNumber string

	stmt, err := DB.Prepare(`SELECT "phone_number" FROM "user" WHERE "user_id" = $1`)
	if err != nil {
		return phoneNumber, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID)

	err = row.Scan(&phoneNumber)
	if err != nil {
		return phoneNumber, fmt.Errorf("error scanning row: %w", err)
	}

	return phoneNumber, nil
}

func GetUserById(id int) (models.User, error) {
	var user models.User

	stmt, err := DB.Prepare(`SELECT user_id, username, password, is_admin, phone_number, first_name, last_name FROM "user" WHERE "user_id" = $1`)
	if err != nil {
		return user, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	err = row.Scan(&user.UserID, &user.Username, &user.Password, &user.IsAdmin, &user.PhoneNumber, &user.FirstName, &user.LastName)
	if err != nil {
		return user, fmt.Errorf("error scanning row: %w", err)
	}

	return user, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User

	stmt, err := DB.Prepare(`SELECT user_id, username, password, is_admin, phone_number, first_name, last_name FROM "user" WHERE "username" = $1`)
	if err != nil {
		return user, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(username)

	err = row.Scan(&user.UserID, &user.Username, &user.Password, &user.IsAdmin, &user.PhoneNumber, &user.FirstName, &user.LastName)
	if err != nil {
		return user, fmt.Errorf("error scanning row: %w", err)
	}

	return user, nil
}

func GetVendingTypes() ([]models.VendingType, error) {
	var vendingTypes []models.VendingType

	rows, err := DB.Query(`SELECT vending_type_id, machine_type FROM "vending_type"`)
	if err != nil {
		return vendingTypes, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var vt models.VendingType
		err := rows.Scan(&vt.VendingTypeID, &vt.MachineType)
		if err != nil {
			return vendingTypes, fmt.Errorf("error scanning row: %w", err)
		}
		vendingTypes = append(vendingTypes, vt)
	}

	if err := rows.Err(); err != nil {
		return vendingTypes, fmt.Errorf("error iterating rows: %w", err)
	}

	return vendingTypes, nil
}

func GetVendingLocations() ([]models.VendingLocation, error) {
	var vendingLocations []models.VendingLocation

	rows, err := DB.Query(`SELECT vending_location_id, location_type FROM "vending_location"`)
	if err != nil {
		return vendingLocations, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var vl models.VendingLocation
		err := rows.Scan(&vl.VendingLocationID, &vl.LocationType)
		if err != nil {
			return vendingLocations, fmt.Errorf("error scanning row: %w", err)
		}
		vendingLocations = append(vendingLocations, vl)
	}

	if err := rows.Err(); err != nil {
		return vendingLocations, fmt.Errorf("error iterating rows: %w", err)
	}

	return vendingLocations, nil
}

func GetLeadList(params types.GetLeadsParams) ([]types.LeadList, int, error) {
	var leads []types.LeadList

	query := `SELECT l.lead_id, l.first_name, l.last_name, l.phone_number, 
		l.created_at, l.rent, l.foot_traffic, l.foot_traffic_type, 
		vt.machine_type, vl.location_type, lm.language, l.vending_type_id, l.vending_location_id,
		COUNT(*) OVER() AS total_rows
		FROM lead AS l
		JOIN vending_type AS vt ON vt.vending_type_id = l.vending_type_id
		JOIN vending_location AS vl ON vl.vending_location_id = l.vending_location_id
		JOIN lead_marketing AS lm ON lm.lead_id = l.lead_id
		WHERE (vt.vending_type_id = $1 OR $1 IS NULL) 
		AND (vl.vending_location_id = $2 OR $2 IS NULL)
		LIMIT $3
		OFFSET $4`

	var offset int

	// Handle pagination
	if params.PageNum != nil {
		pageNum, err := strconv.Atoi(*params.PageNum)
		if err != nil {
			return nil, 0, fmt.Errorf("could not convert page num: %w", err)
		}
		offset = (pageNum - 1) * int(constants.LeadsPerPage)
	}

	rows, err := DB.Query(query, params.VendingType, params.LocationType, constants.LeadsPerPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var totalRows int
	for rows.Next() {
		var lead types.LeadList
		var createdAt time.Time

		var rent, footTraffic, footTrafficType sql.NullString

		err := rows.Scan(&lead.LeadID,
			&lead.FirstName,
			&lead.LastName,
			&lead.PhoneNumber,
			&createdAt,
			&rent,
			&footTraffic,
			&footTrafficType,
			&lead.MachineType,
			&lead.LocationType,
			&lead.Language,
			&lead.VendingTypeID,
			&lead.VendingLocationID,
			&totalRows)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning row: %w", err)
		}
		lead.CreatedAt = createdAt.Unix()

		if rent.Valid {
			lead.Rent = rent.String
		}
		if footTraffic.Valid {
			lead.FootTraffic = footTraffic.String
		}
		if footTrafficType.Valid {
			lead.FootTrafficType = footTrafficType.String
		}

		leads = append(leads, lead)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
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
	l.message
	FROM lead l
	JOIN vending_type vt ON l.vending_type_id = vt.vending_type_id
	JOIN vending_location vl ON l.vending_location_id = vl.vending_location_id
	JOIN lead_marketing lm ON l.lead_id = lm.lead_id
	WHERE l.lead_id = $1`

	var leadDetails types.LeadDetails

	row := DB.QueryRow(query, leadID)

	var adCampaign, medium, source, referrer, landingPage, ip, keyword, channel, language sql.NullString
	var vendingType, vendingLocation, message sql.NullString

	err := row.Scan(
		&leadDetails.LeadID,
		&leadDetails.FirstName,
		&leadDetails.LastName,
		&leadDetails.PhoneNumber,
		&vendingType,
		&vendingLocation,
		&adCampaign,
		&medium,
		&source,
		&referrer,
		&landingPage,
		&ip,
		&keyword,
		&channel,
		&language,
		&message,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return leadDetails, fmt.Errorf("no lead found with ID %s", leadID)
		}
		return leadDetails, fmt.Errorf("error scanning row: %w", err)
	}

	leadDetails.VendingType = vendingType.String
	leadDetails.VendingLocation = vendingLocation.String
	leadDetails.CampaignName = adCampaign.String
	leadDetails.Medium = medium.String
	leadDetails.Source = source.String
	leadDetails.Referrer = referrer.String
	leadDetails.LandingPage = landingPage.String
	leadDetails.IP = ip.String
	leadDetails.Keyword = keyword.String
	leadDetails.Channel = channel.String
	leadDetails.Language = language.String
	leadDetails.Message = message.String

	return leadDetails, nil
}

func GetLeadIDFromPhoneNumber(from string) (int, error) {
	var leadId int

	stmt, err := DB.Prepare(`SELECT "lead_id" FROM "lead" WHERE "phone_number" = $1`)
	if err != nil {
		return leadId, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(from)
	err = row.Scan(&leadId)
	if err != nil {
		return leadId, fmt.Errorf("error scanning row: %w", err)
	}

	return leadId, nil
}

func GetLeadIDFromIncomingTextMessage(from string) (int, error) {
	var leadId int

	stmt, err := DB.Prepare(`SELECT l.lead_id FROM "lead" AS l WHERE l.phone_number = $1`)
	if err != nil {
		return leadId, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(from)

	var forwardingLeadID sql.NullInt64
	err = row.Scan(&forwardingLeadID)
	if err != nil && err != sql.ErrNoRows {
		return leadId, err
	}

	if forwardingLeadID.Valid {
		leadId = int(forwardingLeadID.Int64)
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

		message.DateCreated = utils.FormatTimestamp(dateCreated.Unix())
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
		SET first_name = COALESCE($2, first_name), 
		    last_name = COALESCE($3, last_name), 
		    phone_number = COALESCE($4, phone_number), 
		    vending_type_id = COALESCE($5, vending_type_id), 
		    vending_location_id = COALESCE($6, vending_location_id)
		WHERE lead_id = $1
	`

	args := []interface{}{}
	if form.LeadID != nil {
		args = append(args, *form.LeadID)
	} else {
		return fmt.Errorf("lead_id cannot be nil")
	}

	if form.FirstName != nil {
		args = append(args, *form.FirstName)
	} else {
		args = append(args, nil)
	}
	if form.LastName != nil {
		args = append(args, *form.LastName)
	} else {
		args = append(args, nil)
	}
	if form.PhoneNumber != nil {
		args = append(args, *form.PhoneNumber)
	} else {
		args = append(args, nil)
	}
	if form.VendingType != nil {
		args = append(args, *form.VendingType)
	} else {
		args = append(args, nil)
	}
	if form.VendingLocation != nil {
		args = append(args, *form.VendingLocation)
	} else {
		args = append(args, nil)
	}

	_, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update lead: %v", err)
	}

	return nil
}

func UpdateLeadMarketing(form types.UpdateLeadMarketingForm) error {
	query := `
		UPDATE lead_marketing
		SET ad_campaign = COALESCE($2, ad_campaign), 
		    medium = COALESCE($3, medium), 
		    source = COALESCE($4, source), 
		    referrer = COALESCE($5, referrer), 
		    landing_page = COALESCE($6, landing_page),
		    ip = COALESCE($7, ip), 
		    keyword = COALESCE($8, keyword), 
		    channel = COALESCE($9, channel), 
		    language = COALESCE($10, language)
		WHERE lead_id = $1
	`

	args := []interface{}{}
	if form.LeadID != nil {
		args = append(args, *form.LeadID)
	} else {
		return fmt.Errorf("lead_id cannot be nil")
	}

	if form.CampaignName != nil {
		args = append(args, *form.CampaignName)
	} else {
		args = append(args, nil)
	}
	if form.Medium != nil {
		args = append(args, *form.Medium)
	} else {
		args = append(args, nil)
	}
	if form.Source != nil {
		args = append(args, *form.Source)
	} else {
		args = append(args, nil)
	}
	if form.Referrer != nil {
		args = append(args, *form.Referrer)
	} else {
		args = append(args, nil)
	}
	if form.LandingPage != nil {
		args = append(args, *form.LandingPage)
	} else {
		args = append(args, nil)
	}
	if form.IP != nil {
		args = append(args, *form.IP)
	} else {
		args = append(args, nil)
	}
	if form.Keyword != nil {
		args = append(args, *form.Keyword)
	} else {
		args = append(args, nil)
	}
	if form.Channel != nil {
		args = append(args, *form.Channel)
	} else {
		args = append(args, nil)
	}
	if form.Language != nil {
		args = append(args, *form.Language)
	} else {
		args = append(args, nil)
	}

	_, err := DB.Exec(query, args...)
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

	stmt, err := DB.Prepare(`SELECT phone_call_id, external_id, user_id, lead_id, call_duration, date_created, call_from, call_to, is_inbound, recording_url, status FROM phone_call WHERE external_id = $1`)
	if err != nil {
		return phoneCall, err
	}
	defer stmt.Close()

	var leadID, callDuration sql.NullInt64
	var recordingUrl sql.NullString

	row := stmt.QueryRow(sid)

	err = row.Scan(
		&phoneCall.PhoneCallID,
		&phoneCall.ExternalID,
		&phoneCall.UserID,
		&leadID,
		&callDuration,
		&phoneCall.DateCreated,
		&phoneCall.CallFrom,
		&phoneCall.CallTo,
		&phoneCall.IsInbound,
		&recordingUrl,
		&phoneCall.Status,
	)
	if err != nil {
		return phoneCall, err
	}

	if leadID.Valid {
		phoneCall.LeadID = int(leadID.Int64)
	}

	if callDuration.Valid {
		phoneCall.CallDuration = int(callDuration.Int64)
	}

	if recordingUrl.Valid {
		phoneCall.RecordingURL = recordingUrl.String
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

	var callDuration sql.NullInt32
	var recordingURL sql.NullString
	var status sql.NullString

	if phoneCall.CallDuration != 0 {
		callDuration = sql.NullInt32{Int32: int32(phoneCall.CallDuration), Valid: true}
	} else {
		callDuration = sql.NullInt32{Valid: false}
	}

	if phoneCall.RecordingURL != "" {
		recordingURL = sql.NullString{String: phoneCall.RecordingURL, Valid: true}
	} else {
		recordingURL = sql.NullString{Valid: false}
	}

	if phoneCall.Status != "" {
		status = sql.NullString{String: phoneCall.Status, Valid: true}
	} else {
		status = sql.NullString{Valid: false}
	}

	_, err := DB.Exec(
		query,
		phoneCall.UserID,
		phoneCall.LeadID,
		callDuration,
		phoneCall.DateCreated,
		phoneCall.CallFrom,
		phoneCall.CallTo,
		phoneCall.IsInbound,
		recordingURL,
		status,
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
        SELECT session_id, user_id, csrf_secret, external_id, date_created, date_expires
        FROM sessions
        WHERE csrf_secret = $1
    `
	row := DB.QueryRow(sqlStatement, userKey)

	var dateCreated, dateExpires time.Time
	var userID sql.NullInt32

	err := row.Scan(
		&session.SessionID,
		&userID,
		&session.CSRFSecret,
		&session.ExternalID,
		&dateCreated,
		&dateExpires,
	)
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
        INSERT INTO sessions (csrf_secret, external_id, date_created, date_expires)
        VALUES ($1, $2, to_timestamp($3), to_timestamp($4))
    `

	_, err := DB.Exec(sqlStatement,
		session.CSRFSecret,
		session.ExternalID,
		session.DateCreated,
		session.DateExpires,
	)

	if err != nil {
		return err
	}

	return nil
}

func UpdateSession(session models.Session) error {
	sqlStatement := `
        UPDATE sessions
        SET external_id = $1,
			user_id = $2
        WHERE csrf_secret = $3
    `

	_, err := DB.Exec(sqlStatement,
		session.ExternalID,
		session.UserID,
		session.CSRFSecret,
	)

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

func CreateLeadImage(img models.LeadImage) error {
	stmt, err := DB.Prepare(`
		INSERT INTO lead_image (src, lead_id, date_added, added_by_user_id)
		VALUES ($1, $2, to_timestamp($3), $4)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var leadID sql.NullInt64
	if img.LeadID != 0 {
		leadID = sql.NullInt64{Int64: int64(img.LeadID), Valid: true}
	}

	_, err = stmt.Exec(img.Src, leadID, img.DateAdded, img.AddedByUserID)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CreateLeadNote(note models.LeadNote) error {
	stmt, err := DB.Prepare(`
		INSERT INTO lead_note (note, lead_id, date_added, added_by_user_id)
		VALUES ($1, $2, to_timestamp($3), $4)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var leadID sql.NullInt64
	if note.LeadID != 0 {
		leadID = sql.NullInt64{Int64: int64(note.LeadID), Valid: true}
	}

	_, err = stmt.Exec(note.Note, leadID, note.DateAdded, note.AddedByUserID)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetLeadNotesByLeadID(leadId int) ([]types.FrontendNote, error) {
	var notes []types.FrontendNote

	query := `SELECT u.username,
	n.note,
	n.date_added
	FROM "lead_note" AS n
	JOIN "user" AS u ON u.user_id = n.added_by_user_id
	WHERE n.lead_id = $1
	ORDER BY n.date_added DESC`

	rows, err := DB.Query(query, leadId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return notes, err
	}
	defer rows.Close()

	for rows.Next() {
		var dateAdded time.Time

		var note types.FrontendNote
		err := rows.Scan(
			&note.UserName,
			&note.Note,
			&dateAdded,
		)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return notes, err
		}

		note.DateAdded = utils.FormatTimestamp(dateAdded.Unix())
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return notes, err
	}

	return notes, nil
}

func GetLeadImagesByLeadID(leadId int) ([]models.LeadImage, error) {
	var images []models.LeadImage

	query := fmt.Sprintf(`SELECT '%s' || i.src AS url FROM "lead_image" AS i WHERE i.lead_id = $1;`, constants.AWSS3LiveImagesPath)

	rows, err := DB.Query(query, leadId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return images, err
	}
	defer rows.Close()

	for rows.Next() {

		var image models.LeadImage
		err := rows.Scan(
			&image.Src,
		)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return images, err
		}

		images = append(images, image)
	}

	if err = rows.Err(); err != nil {
		return images, err
	}

	return images, nil
}

func CreateBusiness(form types.BusinessForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO business (name, is_active, date_created, website, industry, google_business_profile)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var dateCreated = time.Now().Unix()

	// Handle NULL values for optional fields
	var name sql.NullString
	var website sql.NullString
	var industry sql.NullString
	var googleBusinessProfile sql.NullString
	var isActive sql.NullBool

	if form.Name != nil {
		name = sql.NullString{String: *form.Name, Valid: true}
	}
	if form.Website != nil {
		website = sql.NullString{String: *form.Website, Valid: true}
	}
	if form.Industry != nil {
		industry = sql.NullString{String: *form.Industry, Valid: true}
	}
	if form.GoogleBusinessProfile != nil {
		googleBusinessProfile = sql.NullString{String: *form.GoogleBusinessProfile, Valid: true}
	}
	if form.IsActive != nil {
		isActive = sql.NullBool{Bool: *form.IsActive, Valid: true}
	}

	_, err = stmt.Exec(
		name,
		isActive,
		dateCreated,
		website,
		industry,
		googleBusinessProfile,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CreateBusinessContact(businessId int, form types.BusinessContactForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO business_contact (
			first_name, 
			last_name, 
			phone, 
			email, 
			preferred_contact_method, 
			preferred_contact_time, 
			business_id, 
			business_position, 
			is_primary_contact
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var preferredContactMethod, preferredContactTime, businessPosition sql.NullString

	if form.PreferredContactMethod != nil {
		preferredContactMethod = sql.NullString{String: *form.PreferredContactMethod, Valid: true}
	}
	if form.PreferredContactTime != nil {
		preferredContactTime = sql.NullString{String: *form.PreferredContactTime, Valid: true}
	}
	if form.BusinessPosition != nil {
		businessPosition = sql.NullString{String: *form.BusinessPosition, Valid: true}
	}

	_, err = stmt.Exec(
		form.FirstName,
		form.LastName,
		form.Phone,
		form.Email,
		preferredContactMethod,
		preferredContactTime,
		businessId,
		businessPosition,
		form.IsPrimaryContact,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CreateLocation(businessId int, form types.LocationForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO location (location_type_id, business_id, name, longitude, latitude, street_address_line_one, street_address_line_two, city_id, zip_code, state, opening, closing)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, to_timestamp($12), to_timestamp($13))
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	// Handle NULL values for optional fields
	var longitude, latitude, streetAddressLineTwo, opening, closing sql.NullString

	if form.Longitude != nil {
		longitude = sql.NullString{String: *form.Longitude, Valid: true}
	}
	if form.Latitude != nil {
		latitude = sql.NullString{String: *form.Latitude, Valid: true}
	}
	if form.StreetAddressLineTwo != nil {
		streetAddressLineTwo = sql.NullString{String: *form.StreetAddressLineTwo, Valid: true}
	}
	if form.Opening != nil {
		opening = sql.NullString{String: *form.Opening, Valid: true}
	}
	if form.Closing != nil {
		closing = sql.NullString{String: *form.Closing, Valid: true}
	}

	_, err = stmt.Exec(
		form.LocationTypeID,
		businessId,
		form.Name,
		longitude,
		latitude,
		form.StreetAddressLineOne,
		streetAddressLineTwo,
		form.CityID,
		form.ZipCode,
		form.State,
		opening,
		closing,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CreateMachine(form types.MachineForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO machine (
			vending_type_id, 
			machine_status_id, 
			location_id, 
			vendor_id,
			year, 
			make, 
			model, 
			purchase_price, 
			purchase_date, 
			card_reader_serial_number, 
			columns_qty, 
			rows_qty, 
			total_slots
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, to_timestamp($9), $10, $11, $12, $13)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	// Handle NULL values for optional fields
	var make, model, cardReaderSerialNumber sql.NullString
	var purchasePrice sql.NullFloat64
	var year, vendingTypeID, locationID, columnsQty, rowsQty, totalSlots, machineStatusID, vendorID sql.NullInt64
	var purchaseDate sql.NullInt64

	if form.Make != nil {
		make = sql.NullString{String: *form.Make, Valid: true}
	}
	if form.Model != nil {
		model = sql.NullString{String: *form.Model, Valid: true}
	}
	if form.CardReaderSerialNumber != nil {
		cardReaderSerialNumber = sql.NullString{String: *form.CardReaderSerialNumber, Valid: true}
	}
	if form.PurchasePrice != nil {
		purchasePrice = sql.NullFloat64{Float64: *form.PurchasePrice, Valid: true}
	}
	if form.Year != nil {
		year = sql.NullInt64{Int64: int64(*form.Year), Valid: true}
	}
	if form.VendingTypeID != nil {
		vendingTypeID = sql.NullInt64{Int64: int64(*form.VendingTypeID), Valid: true}
	}
	if form.LocationID != nil {
		locationID = sql.NullInt64{Int64: int64(*form.LocationID), Valid: true}
	}
	if form.ColumnsQty != nil {
		columnsQty = sql.NullInt64{Int64: int64(*form.ColumnsQty), Valid: true}
	}
	if form.RowsQty != nil {
		rowsQty = sql.NullInt64{Int64: int64(*form.RowsQty), Valid: true}
	}
	if form.TotalSlots != nil {
		totalSlots = sql.NullInt64{Int64: int64(*form.TotalSlots), Valid: true}
	}
	if form.MachineStatusID != nil {
		machineStatusID = sql.NullInt64{Int64: int64(*form.MachineStatusID), Valid: true}
	}
	if form.VendorID != nil {
		vendorID = sql.NullInt64{Int64: int64(*form.VendorID), Valid: true}
	}
	if form.PurchaseDate != nil {
		purchaseDate = sql.NullInt64{Int64: *form.PurchaseDate, Valid: true}
	}

	_, err = stmt.Exec(
		vendingTypeID,
		machineStatusID,
		locationID,
		vendorID,
		year,
		make,
		model,
		purchasePrice,
		purchaseDate,
		cardReaderSerialNumber,
		columnsQty,
		rowsQty,
		totalSlots,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateBusinessContact(businessId int, businessContactId int, form types.BusinessContactForm) error {
	stmt, err := DB.Prepare(`
		UPDATE business_contact
		SET first_name = COALESCE($2, first_name),
		    last_name = COALESCE($3, last_name),
		    phone = COALESCE($4, phone),
		    email = COALESCE($5, email),
		    preferred_contact_method = COALESCE($6, preferred_contact_method),
		    preferred_contact_time = COALESCE($7, preferred_contact_time),
		    business_id = COALESCE($8, business_id),
		    business_position = COALESCE($9, business_position),
		    is_primary_contact = COALESCE($10, is_primary_contact)
		WHERE business_contact_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var firstName, lastName, phone, email, preferredContactMethod, preferredContactTime, businessPosition sql.NullString
	var isPrimaryContact sql.NullBool

	if form.FirstName != nil {
		firstName = sql.NullString{String: *form.FirstName, Valid: true}
	}
	if form.LastName != nil {
		lastName = sql.NullString{String: *form.LastName, Valid: true}
	}
	if form.Phone != nil {
		phone = sql.NullString{String: *form.Phone, Valid: true}
	}
	if form.Email != nil {
		email = sql.NullString{String: *form.Email, Valid: true}
	}
	if form.PreferredContactMethod != nil {
		preferredContactMethod = sql.NullString{String: *form.PreferredContactMethod, Valid: true}
	}
	if form.PreferredContactTime != nil {
		preferredContactTime = sql.NullString{String: *form.PreferredContactTime, Valid: true}
	}
	if form.BusinessPosition != nil {
		businessPosition = sql.NullString{String: *form.BusinessPosition, Valid: true}
	}
	isPrimaryContact = sql.NullBool{Bool: *form.IsPrimaryContact, Valid: true}

	_, err = stmt.Exec(
		businessContactId,
		firstName,
		lastName,
		phone,
		email,
		preferredContactMethod,
		preferredContactTime,
		businessId,
		businessPosition,
		isPrimaryContact,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateBusiness(businessId int, form types.BusinessForm) error {
	stmt, err := DB.Prepare(`
		UPDATE business
		SET name = COALESCE($2, name),
		    website = COALESCE($3, website),
		    industry = COALESCE($4, industry),
		    is_active = COALESCE($5, is_active),
		    google_business_profile = COALESCE($6, google_business_profile)
		WHERE business_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var name, website, industry, googleBusinessProfile sql.NullString
	var isActive sql.NullBool

	if form.Name != nil {
		name = sql.NullString{String: *form.Name, Valid: true}
	}
	if form.Website != nil {
		website = sql.NullString{String: *form.Website, Valid: true}
	}
	if form.Industry != nil {
		industry = sql.NullString{String: *form.Industry, Valid: true}
	}
	if form.GoogleBusinessProfile != nil {
		googleBusinessProfile = sql.NullString{String: *form.GoogleBusinessProfile, Valid: true}
	}
	if form.IsActive != nil {
		isActive = sql.NullBool{Bool: *form.IsActive, Valid: true}
	}

	_, err = stmt.Exec(
		businessId,
		name,
		website,
		industry,
		isActive,
		googleBusinessProfile,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateMachine(machineId int, form types.MachineForm) error {
	stmt, err := DB.Prepare(`
		UPDATE machine
		SET vending_type_id = COALESCE($2, vending_type_id),
		    year = COALESCE($3, year),
		    make = COALESCE($4, make),
		    model = COALESCE($5, model),
		    purchase_price = COALESCE($6, purchase_price),
		    purchase_date = COALESCE($7, purchase_date),
		    card_reader_serial_number = COALESCE($8, card_reader_serial_number),
		    location_id = COALESCE($9, location_id),
		    columns_qty = COALESCE($10, columns_qty),
		    rows_qty = COALESCE($11, rows_qty),
		    total_slots = COALESCE($12, total_slots),
		    machine_status_id = COALESCE($13, machine_status_id),
		    vendor_id = COALESCE($14, vendor_id)
		WHERE machine_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var purchaseDate sql.NullInt64

	if form.PurchaseDate != nil {
		purchaseDate = sql.NullInt64{Int64: *form.PurchaseDate, Valid: true}
	}

	_, err = stmt.Exec(
		machineId,
		form.VendingTypeID,
		form.Year,
		form.Make,
		form.Model,
		form.PurchasePrice,
		purchaseDate,
		form.CardReaderSerialNumber,
		form.LocationID,
		form.ColumnsQty,
		form.RowsQty,
		form.TotalSlots,
		form.MachineStatusID,
		form.VendorID,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateLocation(businessId int, locationId int, form types.LocationForm) error {
	stmt, err := DB.Prepare(`
		UPDATE location
		SET location_type_id = COALESCE($2, location_type_id),
		    business_id = COALESCE($3, business_id),
		    name = COALESCE($4, name),
		    longitude = COALESCE($5, longitude),
		    latitude = COALESCE($6, latitude),
		    street_address_line_one = COALESCE($7, street_address_line_one),
		    street_address_line_two = COALESCE($8, street_address_line_two),
		    city_id = COALESCE($9, city_id),
		    zip_code = COALESCE($10, zip_code),
		    state = COALESCE($11, state),
		    opening = COALESCE($12, opening),
		    closing = COALESCE($13, closing)
		WHERE location_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var longitude, latitude, streetAddressLineTwo, opening, closing sql.NullString

	if form.Longitude != nil {
		longitude = sql.NullString{String: *form.Longitude, Valid: true}
	}
	if form.Latitude != nil {
		latitude = sql.NullString{String: *form.Latitude, Valid: true}
	}
	if form.StreetAddressLineTwo != nil {
		streetAddressLineTwo = sql.NullString{String: *form.StreetAddressLineTwo, Valid: true}
	}
	if form.Opening != nil {
		opening = sql.NullString{String: *form.Opening, Valid: true}
	}
	if form.Closing != nil {
		closing = sql.NullString{String: *form.Closing, Valid: true}
	}

	_, err = stmt.Exec(
		locationId,
		form.LocationTypeID,
		businessId,
		form.Name,
		longitude,
		latitude,
		form.StreetAddressLineOne,
		streetAddressLineTwo,
		form.CityID,
		form.ZipCode,
		form.State,
		opening,
		closing,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetBusinessList(pageNum int) ([]models.Business, int, error) {
	var businesses []models.Business

	rows, err := DB.Query(`SELECT business_id, name, is_active, date_created, website, industry, google_business_profile
	FROM "business"
	ORDER BY date_created DESC
	LIMIT $1
	OFFSET $2;`, constants.LeadsPerPage, pageNum)
	if err != nil {
		return businesses, 0, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var business models.Business
		var dateCreated time.Time
		var website, industry, googleBusinessProfile sql.NullString

		err := rows.Scan(
			&business.BusinessID,
			&business.Name,
			&business.IsActive,
			&dateCreated,
			&website,
			&industry,
			&googleBusinessProfile,
		)
		if err != nil {
			return businesses, 0, fmt.Errorf("error scanning row: %w", err)
		}

		// Handle nullable fields
		if website.Valid {
			business.Website = website.String
		}
		if industry.Valid {
			business.Industry = industry.String
		}
		if googleBusinessProfile.Valid {
			business.GoogleBusinessProfile = googleBusinessProfile.String
		}

		business.DateCreated = dateCreated.Unix()
		businesses = append(businesses, business)
	}

	if err := rows.Err(); err != nil {
		return businesses, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return businesses, len(businesses), nil
}

func GetMachineList(pageNum int) ([]types.MachineList, int, error) {
	var machines []types.MachineList
	var totalRows int

	rows, err := DB.Query(`SELECT CONCAT(m.year, ' ', m.year, ' ', m.model) AS machine_name,
	m.card_reader_serial_number, s.status, l.name, m.purchase_date, COUNT(*) OVER() AS total_rows
	FROM "machine" AS m
	JOIN machine_status AS s ON s.MachineStatusID = m.MachineStatusID
	JOIN location AS l ON l.LocationID = m.LocationID
	ORDER BY m.purchase_date DESC
	LIMIT $1
	OFFSET $2;`, constants.LeadsPerPage, pageNum)
	if err != nil {
		return machines, totalRows, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var machine types.MachineList
		var dateCreated time.Time
		var location, cardReaderSerialNumber sql.NullString

		err := rows.Scan(
			&machine.MachineName,
			&cardReaderSerialNumber,
			&machine.MachineStatus,
			&location,
			&dateCreated,
			&totalRows,
		)
		if err != nil {
			return machines, totalRows, fmt.Errorf("error scanning row: %w", err)
		}

		// Handle nullable fields
		if location.Valid {
			machine.Location = location.String
		}
		if cardReaderSerialNumber.Valid {
			machine.CardReaderSerialNumber = cardReaderSerialNumber.String
		}

		machine.PurchaseDate = dateCreated.Unix()
		machines = append(machines, machine)
	}

	if err := rows.Err(); err != nil {
		return machines, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return machines, totalRows, nil
}

func GetLocations() ([]models.Location, error) {
	var locations []models.Location

	rows, err := DB.Query(`
		SELECT location_id, location_type_id, business_id, name, longitude, latitude, street_address_line_one, street_address_line_two, city_id, zip_code, state, opening, closing 
		FROM "location"
	`)
	if err != nil {
		return locations, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var loc models.Location
		var longitude, latitude, streetAddressLineOne, streetAddressLineTwo, zipCode, state, opening, closing sql.NullString
		var cityID sql.NullInt64

		err := rows.Scan(
			&loc.LocationID,
			&loc.LocationTypeID,
			&loc.BusinessID,
			&loc.Name,
			&longitude,
			&latitude,
			&streetAddressLineOne,
			&streetAddressLineTwo,
			&cityID,
			&zipCode,
			&state,
			&opening,
			&closing,
		)
		if err != nil {
			return locations, fmt.Errorf("error scanning row: %w", err)
		}

		if longitude.Valid {
			loc.Longitude = longitude.String
		}
		if latitude.Valid {
			loc.Latitude = latitude.String
		}
		if streetAddressLineOne.Valid {
			loc.StreetAdressLineOne = streetAddressLineOne.String
		}
		if streetAddressLineTwo.Valid {
			loc.StreetAdressLineTwo = streetAddressLineTwo.String
		}
		if cityID.Valid {
			loc.CityID = int(cityID.Int64)
		}
		if zipCode.Valid {
			loc.ZipCode = zipCode.String
		}
		if state.Valid {
			loc.State = state.String
		}
		if opening.Valid {
			loc.Opening = opening.String
		}
		if closing.Valid {
			loc.Closing = closing.String
		}

		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		return locations, fmt.Errorf("error iterating rows: %w", err)
	}

	return locations, nil
}

func GetMachineStatuses() ([]models.MachineStatus, error) {
	var statuses []models.MachineStatus

	rows, err := DB.Query(`
		SELECT machine_status_id, status 
		FROM "machine_status"
	`)
	if err != nil {
		return statuses, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status models.MachineStatus
		var statusText sql.NullString

		err := rows.Scan(
			&status.MachineStatusID,
			&statusText,
		)
		if err != nil {
			return statuses, fmt.Errorf("error scanning row: %w", err)
		}

		if statusText.Valid {
			status.Status = statusText.String
		}

		statuses = append(statuses, status)
	}

	if err := rows.Err(); err != nil {
		return statuses, fmt.Errorf("error iterating rows: %w", err)
	}

	return statuses, nil
}

func GetVendors() ([]models.Vendor, error) {
	var vendors []models.Vendor

	rows, err := DB.Query(`
		SELECT vendor_id, name, first_name, last_name, phone, email, preferred_contact_method, preferred_contact_time, street_address_line_one, street_address_line_two, city_id, zip_code, state, google_business_profile 
		FROM "vendor"
	`)
	if err != nil {
		return vendors, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var vendor models.Vendor
		var name, firstName, lastName, phone, email, preferredContactMethod, preferredContactTime, streetAddressLineOne, streetAddressLineTwo, zipCode, state, googleBusinessProfile sql.NullString

		err := rows.Scan(
			&vendor.VendorID,
			&name,
			&firstName,
			&lastName,
			&phone,
			&email,
			&preferredContactMethod,
			&preferredContactTime,
			&streetAddressLineOne,
			&streetAddressLineTwo,
			&vendor.CityID,
			&zipCode,
			&state,
			&googleBusinessProfile,
		)
		if err != nil {
			return vendors, fmt.Errorf("error scanning row: %w", err)
		}

		if name.Valid {
			vendor.Name = name.String
		}
		if firstName.Valid {
			vendor.FirstName = firstName.String
		}
		if lastName.Valid {
			vendor.LastName = lastName.String
		}
		if phone.Valid {
			vendor.Phone = phone.String
		}
		if email.Valid {
			vendor.Email = email.String
		}
		if preferredContactMethod.Valid {
			vendor.PreferredContactMethod = preferredContactMethod.String
		}
		if preferredContactTime.Valid {
			vendor.PreferredContactTime = preferredContactTime.String
		}
		if streetAddressLineOne.Valid {
			vendor.StreetAdressLineOne = streetAddressLineOne.String
		}
		if streetAddressLineTwo.Valid {
			vendor.StreetAdressLineTwo = streetAddressLineTwo.String
		}
		if zipCode.Valid {
			vendor.ZipCode = zipCode.String
		}
		if state.Valid {
			vendor.State = state.String
		}
		if googleBusinessProfile.Valid {
			vendor.GoogleBusinessProfile = googleBusinessProfile.String
		}

		vendors = append(vendors, vendor)
	}

	if err := rows.Err(); err != nil {
		return vendors, fmt.Errorf("error iterating rows: %w", err)
	}

	return vendors, nil
}
