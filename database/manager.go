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

func GetDashboardStats() (types.DashboardStats, error) {
	var counts types.DashboardStats

	query := `
        SELECT 
            (SELECT COUNT(1) FROM lead) AS leads,
            (SELECT COUNT(1) FROM business) AS businesses,
            (SELECT COUNT(1) FROM vendor) AS vendors,
            (SELECT COUNT(1) FROM supplier) AS suppliers,
            (SELECT COUNT(1) FROM machine) AS machines;
    `

	row := DB.QueryRow(query)
	err := row.Scan(
		&counts.Leads,
		&counts.Businesses,
		&counts.Vendors,
		&counts.Suppliers,
		&counts.Machines,
	)
	if err != nil {
		return counts, fmt.Errorf("error scanning row: %w", err)
	}

	return counts, nil
}

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

	// Convert fields using utility functions
	createdAt := time.Now().Unix()
	rent := utils.CreateNullString(quoteForm.Rent)
	footTraffic := utils.CreateNullString(quoteForm.FootTraffic)
	footTrafficType := utils.CreateNullString(quoteForm.FootTrafficType)
	message := utils.CreateNullString(quoteForm.Message)
	vendingTypeID := utils.CreateNullInt(quoteForm.MachineType)
	vendingLocationID := utils.CreateNullInt(quoteForm.LocationType)

	err = leadStmt.QueryRow(
		utils.CreateNullString(quoteForm.FirstName),
		utils.CreateNullString(quoteForm.LastName),
		utils.CreateNullString(quoteForm.PhoneNumber),
		createdAt,
		rent,
		footTraffic,
		footTrafficType,
		vendingTypeID,
		vendingLocationID,
		message,
	).Scan(&leadID)
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

	_, err = marketingStmt.Exec(
		leadID,
		utils.CreateNullString(quoteForm.Source),
		utils.CreateNullString(quoteForm.Medium),
		utils.CreateNullString(quoteForm.Channel),
		utils.CreateNullString(quoteForm.LandingPage),
		utils.CreateNullString(quoteForm.Keyword),
		utils.CreateNullString(quoteForm.Referrer),
		utils.CreateNullString(quoteForm.ClickID),
		utils.CreateNullInt64(quoteForm.CampaignID),
		utils.CreateNullString(quoteForm.AdCampaign),
		utils.CreateNullInt64(quoteForm.AdGroupID),
		utils.CreateNullString(quoteForm.AdGroupName),
		utils.CreateNullInt64(quoteForm.AdSetID),
		utils.CreateNullString(quoteForm.AdSetName),
		utils.CreateNullInt64(quoteForm.AdID),
		utils.CreateNullInt64(quoteForm.AdHeadline),
		utils.CreateNullString(quoteForm.Language),
		utils.CreateNullString(quoteForm.UserAgent),
		utils.CreateNullString(quoteForm.ButtonClicked),
		utils.CreateNullString(quoteForm.IP),
		utils.CreateNullString(quoteForm.ExternalID),
		utils.CreateNullString(quoteForm.GoogleClientID),
		utils.CreateNullString(quoteForm.CSRFSecret),
		utils.CreateNullString(quoteForm.FacebookClickID),
		utils.CreateNullString(quoteForm.FacebookClientID),
		utils.CreateNullString(quoteForm.Longitude),
		utils.CreateNullString(quoteForm.Latitude),
	)
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
	if form.LeadID == nil {
		return fmt.Errorf("lead_id cannot be nil")
	}

	query := `
		UPDATE lead
		SET first_name = COALESCE($2, first_name), 
		    last_name = COALESCE($3, last_name), 
		    phone_number = COALESCE($4, phone_number), 
		    vending_type_id = COALESCE($5, vending_type_id), 
		    vending_location_id = COALESCE($6, vending_location_id)
		WHERE lead_id = $1
	`

	args := []interface{}{
		*form.LeadID,
		utils.CreateNullString(form.FirstName),
		utils.CreateNullString(form.LastName),
		utils.CreateNullString(form.PhoneNumber),
		utils.CreateNullInt(form.VendingType),
		utils.CreateNullInt(form.VendingLocation),
	}

	_, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update lead: %v", err)
	}

	return nil
}

func UpdateLeadMarketing(form types.UpdateLeadMarketingForm) error {
	if form.LeadID == nil {
		return fmt.Errorf("lead_id cannot be nil")
	}

	query := `
		UPDATE lead_marketing
		SET ad_campaign = $2, 
		    medium = $3, 
		    source = $4, 
		    referrer = $5, 
		    landing_page = $6,
		    ip = $7, 
		    keyword = $8, 
		    channel = $9, 
		    language = $10
		WHERE lead_id = $1
	`

	args := []interface{}{
		*form.LeadID,
		utils.CreateNullString(form.CampaignName),
		utils.CreateNullString(form.Medium),
		utils.CreateNullString(form.Source),
		utils.CreateNullString(form.Referrer),
		utils.CreateNullString(form.LandingPage),
		utils.CreateNullString(form.IP),
		utils.CreateNullString(form.Keyword),
		utils.CreateNullString(form.Channel),
		utils.CreateNullString(form.Language),
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
			call_duration = COALESCE($3, call_duration),
			date_created = $4,
			call_from = $5,
			call_to = $6,
			is_inbound = $7,
			recording_url = COALESCE($8, recording_url),
			status = COALESCE($9, status)
		WHERE external_id = $10`

	args := []interface{}{
		phoneCall.UserID,
		phoneCall.LeadID,
		utils.CreateNullInt(&phoneCall.CallDuration),
		phoneCall.DateCreated,
		phoneCall.CallFrom,
		phoneCall.CallTo,
		phoneCall.IsInbound,
		utils.CreateNullString(&phoneCall.RecordingURL),
		utils.CreateNullString(&phoneCall.Status),
		phoneCall.ExternalID,
	}

	_, err := DB.Exec(query, args...)
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
        SET external_id = COALESCE($1, external_id),
            user_id = COALESCE($2, user_id)
        WHERE csrf_secret = $3
    `

	args := []interface{}{
		utils.CreateNullString(&session.ExternalID),
		utils.CreateNullInt(&session.UserID),
		session.CSRFSecret,
	}

	_, err := DB.Exec(sqlStatement, args...)
	if err != nil {
		return fmt.Errorf("error updating session: %w", err)
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

	leadID := utils.CreateNullInt(&img.LeadID)

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

	leadID := utils.CreateNullInt(&note.LeadID)

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
		INSERT INTO business (name, is_active, website, industry, google_business_profile)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	name := utils.CreateNullString(form.Name)
	isActive := utils.CreateNullBool(form.IsActive)
	website := utils.CreateNullString(form.Website)
	industry := utils.CreateNullString(form.Industry)
	googleBusinessProfile := utils.CreateNullString(form.GoogleBusinessProfile)

	_, err = stmt.Exec(
		name,
		isActive,
		website,
		industry,
		googleBusinessProfile,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CreateBusinessContact(form types.BusinessContactForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO business_contact (
			first_name, 
			last_name, 
			phone, 
			email, 
			preferred_contact_method, 
			preferred_contact_time, 
			business_position
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	preferredContactMethod := utils.CreateNullString(form.PreferredContactMethod)
	preferredContactTime := utils.CreateNullString(form.PreferredContactTime)
	businessPosition := utils.CreateNullString(form.BusinessPosition)

	_, err = stmt.Exec(
		form.FirstName,
		form.LastName,
		form.Phone,
		form.Email,
		preferredContactMethod,
		preferredContactTime,
		businessPosition,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CreateLocation(businessId int, form types.LocationForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO location (vending_location_id, business_id, name, longitude, latitude, street_address_line_one, street_address_line_two, city_id, zip_code, state, opening, closing, date_started)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, to_timestamp($13))
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	longitude := utils.CreateNullFloat64(form.Longitude)
	latitude := utils.CreateNullFloat64(form.Latitude)
	streetAddressLineTwo := utils.CreateNullString(form.StreetAddressLineTwo)
	opening := utils.CreateNullString(form.Opening)
	closing := utils.CreateNullString(form.Closing)

	_, err = stmt.Exec(
		form.VendingLocationID,
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
		form.DateStarted,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CreateVendor(form types.VendorForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO vendor (
			name,
			first_name,
			last_name,
			phone_number,
			email,
			preferred_contact_method,
			preferred_contact_time,
			street_address_line_one,
			street_address_line_two,
			city_id,
			zip_code,
			state,
			google_business_profile
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	name := utils.CreateNullString(form.Name)
	firstName := utils.CreateNullString(form.FirstName)
	lastName := utils.CreateNullString(form.LastName)
	phone := utils.CreateNullString(form.PhoneNumber)
	email := utils.CreateNullString(form.Email)
	preferredContactMethod := utils.CreateNullString(form.PreferredContactMethod)
	preferredContactTime := utils.CreateNullString(form.PreferredContactTime)
	streetAddressLineOne := utils.CreateNullString(form.StreetAddressLineOne)
	streetAddressLineTwo := utils.CreateNullString(form.StreetAddressLineTwo)
	zipCode := utils.CreateNullString(form.ZipCode)
	state := utils.CreateNullString(form.State)
	googleBusinessProfile := utils.CreateNullString(form.GoogleBusinessProfile)
	cityID := utils.CreateNullInt(form.CityID)

	_, err = stmt.Exec(
		name,
		firstName,
		lastName,
		phone,
		email,
		preferredContactMethod,
		preferredContactTime,
		streetAddressLineOne,
		streetAddressLineTwo,
		cityID,
		zipCode,
		state,
		googleBusinessProfile,
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
		    preferred_contact_method = $6,
		    preferred_contact_time = $7,
		    business_id = COALESCE($8, business_id),
		    business_position = $9
		WHERE business_contact_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	firstName := utils.CreateNullString(form.FirstName)
	lastName := utils.CreateNullString(form.LastName)
	phone := utils.CreateNullString(form.Phone)
	email := utils.CreateNullString(form.Email)
	preferredContactMethod := utils.CreateNullString(form.PreferredContactMethod)
	preferredContactTime := utils.CreateNullString(form.PreferredContactTime)
	businessPosition := utils.CreateNullString(form.BusinessPosition)

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
		    website = $3,
		    industry = $4,
		    is_active = COALESCE($5, is_active),
		    google_business_profile = $6
		WHERE business_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	name := utils.CreateNullString(form.Name)
	website := utils.CreateNullString(form.Website)
	industry := utils.CreateNullString(form.Industry)
	googleBusinessProfile := utils.CreateNullString(form.GoogleBusinessProfile)
	isActive := utils.CreateNullBool(form.IsActive)

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
		    year = $3,
		    make = $4,
		    model = $5,
		    purchase_price = $6,
		    purchase_date = to_timestamp($7),
		    card_reader_serial_number = $8,
		    location_id = $9,
		    machine_status_id = COALESCE($10, machine_status_id),
		    vendor_id = COALESCE($11, vendor_id)
		WHERE machine_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	vendingTypeID := utils.CreateNullInt(form.VendingTypeID)
	year := utils.CreateNullInt(form.Year)
	make := utils.CreateNullString(form.Make)
	model := utils.CreateNullString(form.Model)
	purchasePrice := utils.CreateNullFloat64(form.PurchasePrice)
	purchaseDate := utils.CreateNullInt64(form.PurchaseDate)
	cardReaderSerialNumber := utils.CreateNullString(form.CardReaderSerialNumber)
	locationID := utils.CreateNullInt(form.LocationID)
	machineStatusID := utils.CreateNullInt(form.MachineStatusID)
	vendorID := utils.CreateNullInt(form.VendorID)

	_, err = stmt.Exec(
		machineId,
		vendingTypeID,
		year,
		make,
		model,
		purchasePrice,
		purchaseDate,
		cardReaderSerialNumber,
		locationID,
		machineStatusID,
		vendorID,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateLocation(businessId int, locationId int, form types.LocationForm) error {
	stmt, err := DB.Prepare(`
		UPDATE location
		SET vending_location_id = COALESCE($2, vending_location_id),
		    business_id = COALESCE($3, business_id),
		    name = COALESCE($4, name),
		    longitude = $5,
		    latitude = $6,
		    street_address_line_one = COALESCE($7, street_address_line_one),
		    street_address_line_two = $8,
		    city_id = COALESCE($9, city_id),
		    zip_code = COALESCE($10, zip_code),
		    state = COALESCE($11, state),
		    opening = $12,
		    closing = $13,
		    date_started = COALESCE(to_timestamp($14), date_started)
		WHERE location_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	longitude := utils.CreateNullFloat64(form.Longitude)
	latitude := utils.CreateNullFloat64(form.Latitude)
	streetAddressLineTwo := utils.CreateNullString(form.StreetAddressLineTwo)
	opening := utils.CreateNullString(form.Opening)
	closing := utils.CreateNullString(form.Closing)

	_, err = stmt.Exec(
		locationId,
		form.VendingLocationID,
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
		form.DateStarted,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetBusinessList(pageNum int) ([]models.Business, int, error) {
	var businesses []models.Business

	var offset = (pageNum - 1) * int(constants.LeadsPerPage)

	rows, err := DB.Query(`SELECT business_id, name, is_active, website, industry, google_business_profile
	FROM "business"
	LIMIT $1
	OFFSET $2;`, constants.LeadsPerPage, offset)
	if err != nil {
		return businesses, 0, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var business models.Business
		var website, industry, googleBusinessProfile sql.NullString

		err := rows.Scan(
			&business.BusinessID,
			&business.Name,
			&business.IsActive,
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

	var offset = (pageNum - 1) * int(constants.LeadsPerPage)

	rows, err := DB.Query(`SELECT m.machine_id,
	CONCAT(m.year, ' ', m.make, ' ', m.model) AS machine_name,
	m.card_reader_serial_number, s.status, l.name, m.purchase_date, COUNT(*) OVER() AS total_rows
	FROM "machine" AS m
	JOIN machine_status AS s ON s.machine_status_id = m.machine_status_id
	LEFT JOIN location AS l ON l.location_id = m.location_id
	ORDER BY m.purchase_date DESC
	LIMIT $1
	OFFSET $2;`, constants.LeadsPerPage, offset)
	if err != nil {
		return machines, totalRows, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var machine types.MachineList
		var dateCreated time.Time
		var location, cardReaderSerialNumber sql.NullString

		err := rows.Scan(
			&machine.MachineID,
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

		machine.PurchaseDate = utils.FormatDateMMDDYYYY(dateCreated.Unix())
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
		SELECT location_id, vending_location_id, name, longitude, latitude, street_address_line_one, street_address_line_two, city_id, zip_code, state, opening, closing, date_started 
		FROM "location"
	`)
	if err != nil {
		return locations, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var loc models.Location
		var streetAddressLineOne, streetAddressLineTwo, zipCode, state, opening, closing sql.NullString
		var cityID sql.NullInt64
		var longitude, latitude sql.NullFloat64
		var dateStarted time.Time

		err := rows.Scan(
			&loc.LocationID,
			&loc.VendingLocationID,
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
			&dateStarted,
		)
		if err != nil {
			return locations, fmt.Errorf("error scanning row: %w", err)
		}

		if longitude.Valid {
			loc.Longitude = longitude.Float64
		}
		if latitude.Valid {
			loc.Latitude = latitude.Float64
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

		loc.DateStarted = dateStarted.Unix()

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

func GetCities() ([]models.City, error) {
	var cities []models.City

	rows, err := DB.Query(`SELECT "city_id", "name" FROM "city"`)
	if err != nil {
		return cities, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var city models.City
		err := rows.Scan(&city.CityID, &city.Name)
		if err != nil {
			return cities, fmt.Errorf("error scanning row: %w", err)
		}
		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
		return cities, fmt.Errorf("error iterating rows: %w", err)
	}

	return cities, nil
}

func GetVendors() ([]models.Vendor, error) {
	var vendors []models.Vendor

	rows, err := DB.Query(`
		SELECT vendor_id, name, first_name, last_name, phone_number, email, preferred_contact_method, preferred_contact_time, street_address_line_one, street_address_line_two, city_id, zip_code, state, google_business_profile 
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
			vendor.PhoneNumber = phone.String
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
			vendor.StreetAddressLineOne = streetAddressLineOne.String
		}
		if streetAddressLineTwo.Valid {
			vendor.StreetAddressLineTwo = streetAddressLineTwo.String
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

func GetBusinesses() ([]models.Business, error) {
	var businesses []models.Business

	rows, err := DB.Query(`
		SELECT business_id, name, is_active, website, industry, google_business_profile
		FROM business
	`)
	if err != nil {
		return businesses, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var business models.Business
		var website, industry, googleBusinessProfile sql.NullString

		err := rows.Scan(
			&business.BusinessID,
			&business.Name,
			&business.IsActive,
			&website,
			&industry,
			&googleBusinessProfile,
		)
		if err != nil {
			return businesses, fmt.Errorf("error scanning row: %w", err)
		}

		// Check for null values and assign them if valid
		if website.Valid {
			business.Website = website.String
		}
		if industry.Valid {
			business.Industry = industry.String
		}
		if googleBusinessProfile.Valid {
			business.GoogleBusinessProfile = googleBusinessProfile.String
		}

		businesses = append(businesses, business)
	}

	if err := rows.Err(); err != nil {
		return businesses, fmt.Errorf("error iterating rows: %w", err)
	}

	return businesses, nil
}

func GetVendorList(pageNum int) ([]types.VendorList, int, error) {
	var vendors []types.VendorList
	var totalRows int

	var offset = (pageNum - 1) * int(constants.LeadsPerPage)

	rows, err := DB.Query(`
		SELECT 
			v.vendor_id,
			v.name,
			v.first_name,
			v.last_name,
			v.phone_number,
			v.email,
			v.preferred_contact_method,
			v.preferred_contact_time,
			v.street_address_line_one,
			v.street_address_line_two,
			v.city_id,
			v.zip_code,
			v.state,
			v.google_business_profile,
			c.name AS city_name,
			COUNT(*) OVER() AS total_rows
		FROM "vendor" AS v
		JOIN city AS c ON c.city_id = v.city_id
		ORDER BY v.vendor_id DESC
		LIMIT $1
		OFFSET $2;
	`, constants.LeadsPerPage, offset)
	if err != nil {
		return vendors, totalRows, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var vendor types.VendorList
		var googleBusinessProfile sql.NullString
		var streetAddressLineOne, streetAddressLineTwo sql.NullString
		var cityName sql.NullString

		err := rows.Scan(
			&vendor.VendorID,
			&vendor.Name,
			&vendor.FirstName,
			&vendor.LastName,
			&vendor.PhoneNumber,
			&vendor.Email,
			&vendor.PreferredContactMethod,
			&vendor.PreferredContactTime,
			&streetAddressLineOne,
			&streetAddressLineTwo,
			&vendor.CityID,
			&vendor.ZipCode,
			&vendor.State,
			&googleBusinessProfile,
			&cityName,
			&totalRows,
		)
		if err != nil {
			return vendors, totalRows, fmt.Errorf("error scanning row: %w", err)
		}

		// Handle nullable fields
		if streetAddressLineOne.Valid {
			vendor.StreetAddressLineOne = streetAddressLineOne.String
		}
		if streetAddressLineTwo.Valid {
			vendor.StreetAddressLineTwo = streetAddressLineTwo.String
		}
		if cityName.Valid {
			vendor.CityName = cityName.String
		}
		if googleBusinessProfile.Valid {
			vendor.GoogleBusinessProfile = googleBusinessProfile.String
		}

		vendors = append(vendors, vendor)
	}

	if err := rows.Err(); err != nil {
		return vendors, totalRows, fmt.Errorf("error iterating rows: %w", err)
	}

	return vendors, totalRows, nil
}

func UpdateVendor(vendorId int, form types.VendorForm) error {
	stmt, err := DB.Prepare(`
		UPDATE vendor
		SET name = $2,
		    first_name = $3,
		    last_name = $4,
		    phone_number = $5,
		    email = $6,
		    preferred_contact_method = $7,
		    preferred_contact_time = $8,
		    street_address_line_one = $9,
		    street_address_line_two = $10,
		    city_id = $11,
		    zip_code = $12,
		    state = $13,
		    google_business_profile = $14
		WHERE vendor_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	name := utils.CreateNullString(form.Name)
	firstName := utils.CreateNullString(form.FirstName)
	lastName := utils.CreateNullString(form.LastName)
	phone := utils.CreateNullString(form.PhoneNumber)
	email := utils.CreateNullString(form.Email)
	preferredContactMethod := utils.CreateNullString(form.PreferredContactMethod)
	preferredContactTime := utils.CreateNullString(form.PreferredContactTime)
	streetAddressLineOne := utils.CreateNullString(form.StreetAddressLineOne)
	streetAddressLineTwo := utils.CreateNullString(form.StreetAddressLineTwo)
	zipCode := utils.CreateNullString(form.ZipCode)
	state := utils.CreateNullString(form.State)
	googleBusinessProfile := utils.CreateNullString(form.GoogleBusinessProfile)
	cityID := utils.CreateNullInt(form.CityID)

	_, err = stmt.Exec(
		vendorId,
		name,
		firstName,
		lastName,
		phone,
		email,
		preferredContactMethod,
		preferredContactTime,
		streetAddressLineOne,
		streetAddressLineTwo,
		cityID,
		zipCode,
		state,
		googleBusinessProfile,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetSuppliers() ([]models.Supplier, error) {
	var suppliers []models.Supplier

	rows, err := DB.Query(`
		SELECT supplier_id, name, membership_id, membership_cost::NUMERIC, membership_renewal, street_address_line_one, street_address_line_two, city_id, zip_code, state, google_business_profile 
		FROM "supplier"
	`)
	if err != nil {
		return suppliers, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var supplier models.Supplier
		var googleBusinessProfile, membershipRenewal, membershipID, streetAddressLineTwo sql.NullString
		var membershipCost sql.NullFloat64

		err := rows.Scan(
			&supplier.SupplierID,
			&supplier.Name,
			&membershipID,
			&membershipCost,
			&membershipRenewal,
			&supplier.StreetAddressLineOne,
			&streetAddressLineTwo,
			&supplier.CityID,
			&supplier.ZipCode,
			&supplier.State,
			&googleBusinessProfile,
		)
		if err != nil {
			return suppliers, fmt.Errorf("error scanning row: %w", err)
		}

		if googleBusinessProfile.Valid {
			supplier.GoogleBusinessProfile = googleBusinessProfile.String
		}
		if membershipRenewal.Valid {
			supplier.MembershipRenewal = membershipRenewal.String
		}
		if membershipID.Valid {
			supplier.MembershipID = membershipID.String
		}
		if membershipCost.Valid {
			supplier.MembershipCost = membershipCost.Float64
		}
		if streetAddressLineTwo.Valid {
			supplier.StreetAddressLineTwo = streetAddressLineTwo.String
		}

		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		return suppliers, fmt.Errorf("error iterating rows: %w", err)
	}

	return suppliers, nil
}

func GetSupplierList(pageNum int) ([]types.SupplierList, int, error) {
	var suppliers []types.SupplierList
	var totalRows int

	var offset = (pageNum - 1) * int(constants.LeadsPerPage)

	rows, err := DB.Query(`
		SELECT 
			s.supplier_id,
			s.name,
			s.membership_id,
			s.membership_cost::NUMERIC,
			s.membership_renewal,
			s.street_address_line_one,
			s.street_address_line_two,
			s.zip_code,
			s.state,
			s.google_business_profile,
			c.name,
			COUNT(*) OVER() AS total_rows
		FROM "supplier" AS s
		JOIN city AS c ON c.city_id = s.city_id
		ORDER BY s.supplier_id DESC
		LIMIT $1
		OFFSET $2;
	`, constants.LeadsPerPage, offset)
	if err != nil {
		return suppliers, totalRows, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var supplier types.SupplierList
		var googleBusinessProfile sql.NullString
		var streetAddressLineOne, streetAddressLineTwo, membershipRenewal, membershipID sql.NullString
		var cityName sql.NullString
		var membershipCost sql.NullFloat64

		err := rows.Scan(
			&supplier.SupplierID,
			&supplier.Name,
			&membershipID,
			&membershipCost,
			&membershipRenewal,
			&streetAddressLineOne,
			&streetAddressLineTwo,
			&supplier.ZipCode,
			&supplier.State,
			&googleBusinessProfile,
			&cityName,
			&totalRows,
		)
		if err != nil {
			return suppliers, totalRows, fmt.Errorf("error scanning row: %w", err)
		}

		// Handle nullable fields
		if streetAddressLineOne.Valid {
			supplier.StreetAddressLineOne = streetAddressLineOne.String
		}
		if streetAddressLineTwo.Valid {
			supplier.StreetAddressLineTwo = streetAddressLineTwo.String
		}
		if cityName.Valid {
			supplier.City = cityName.String
		}
		if googleBusinessProfile.Valid {
			supplier.GoogleBusinessProfile = googleBusinessProfile.String
		}
		if membershipRenewal.Valid {
			supplier.MembershipRenewal = membershipRenewal.String
		}
		if membershipID.Valid {
			supplier.MembershipID = membershipID.String
		}
		if membershipCost.Valid {
			supplier.MembershipCost = membershipCost.Float64
		}

		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		return suppliers, totalRows, fmt.Errorf("error iterating rows: %w", err)
	}

	return suppliers, totalRows, nil
}

func CreateSupplier(form types.SupplierForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO supplier (
			name,
			membership_id,
			membership_cost,
			membership_renewal,
			street_address_line_one,
			street_address_line_two,
			city_id,
			zip_code,
			state,
			google_business_profile
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	membershipID := utils.CreateNullString(form.MembershipID)
	membershipCost := utils.CreateNullFloat64(form.MembershipCost)
	membershipRenewal := utils.CreateNullString(form.MembershipRenewal)
	streetAddressLineTwo := utils.CreateNullString(form.StreetAddressLineTwo)
	googleBusinessProfile := utils.CreateNullString(form.GoogleBusinessProfile)

	_, err = stmt.Exec(
		form.Name,
		membershipID,
		membershipCost,
		membershipRenewal,
		form.StreetAddressLineOne,
		streetAddressLineTwo,
		form.CityID,
		form.ZipCode,
		form.State,
		googleBusinessProfile,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateSupplier(supplierId int, form types.SupplierForm) error {
	stmt, err := DB.Prepare(`
		UPDATE supplier
		SET name = $2,
		    membership_id = $3,
		    membership_cost = $4,
		    membership_renewal = $5,
		    street_address_line_one = $6,
		    street_address_line_two = $7,
		    city_id = $8,
		    zip_code = $9,
		    state = $10,
		    google_business_profile = $11
		WHERE supplier_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	name := utils.CreateNullString(form.Name)
	membershipID := utils.CreateNullString(form.MembershipID)
	membershipCost := utils.CreateNullFloat64(form.MembershipCost)
	membershipRenewal := utils.CreateNullString(form.MembershipRenewal)
	streetAddressLineOne := utils.CreateNullString(form.StreetAddressLineOne)
	streetAddressLineTwo := utils.CreateNullString(form.StreetAddressLineTwo)
	cityID := utils.CreateNullInt(form.CityID)
	zipCode := utils.CreateNullString(form.ZipCode)
	state := utils.CreateNullString(form.State)
	googleBusinessProfile := utils.CreateNullString(form.GoogleBusinessProfile)

	_, err = stmt.Exec(
		supplierId,
		name,
		membershipID,
		membershipCost,
		membershipRenewal,
		streetAddressLineOne,
		streetAddressLineTwo,
		cityID,
		zipCode,
		state,
		googleBusinessProfile,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func DeleteMachine(id int) error {
	sqlStatement := `
        DELETE FROM machine WHERE machine_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteLocation(id int) error {
	sqlStatement := `
        DELETE FROM location WHERE location_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteBusiness(id int) error {
	sqlStatement := `
        DELETE FROM business WHERE business_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSupplier(id int) error {
	sqlStatement := `
        DELETE FROM supplier WHERE supplier_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteVendor(id int) error {
	sqlStatement := `
        DELETE FROM vendor WHERE vendor_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func GetVendorDetails(vendorID string) (types.VendorDetails, error) {
	query := `SELECT 
		v.vendor_id,
		v.name,
		v.first_name,
		v.last_name,
		v.phone_number,
		v.email,
		v.preferred_contact_method,
		v.preferred_contact_time,
		v.street_address_line_one,
		v.street_address_line_two,
		v.city_id,
		v.zip_code,
		v.state,
		v.google_business_profile
	FROM vendor AS v
	WHERE v.vendor_id = $1`

	var vendorDetails types.VendorDetails

	row := DB.QueryRow(query, vendorID)

	var (
		googleBusinessProfile  sql.NullString
		streetAddressLineOne   sql.NullString
		streetAddressLineTwo   sql.NullString
		preferredContactMethod sql.NullString
		preferredContactTime   sql.NullString
	)

	err := row.Scan(
		&vendorDetails.VendorID,
		&vendorDetails.Name,
		&vendorDetails.FirstName,
		&vendorDetails.LastName,
		&vendorDetails.PhoneNumber,
		&vendorDetails.Email,
		&preferredContactMethod,
		&preferredContactTime,
		&streetAddressLineOne,
		&streetAddressLineTwo,
		&vendorDetails.CityID,
		&vendorDetails.ZipCode,
		&vendorDetails.State,
		&googleBusinessProfile,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return vendorDetails, fmt.Errorf("no vendor found with ID %s", vendorID)
		}
		return vendorDetails, fmt.Errorf("error scanning row: %w", err)
	}

	// Handle nullable fields
	if preferredContactMethod.Valid {
		vendorDetails.PreferredContactMethod = preferredContactMethod.String
	}
	if preferredContactTime.Valid {
		vendorDetails.PreferredContactTime = preferredContactTime.String
	}
	if streetAddressLineOne.Valid {
		vendorDetails.StreetAddressLineOne = streetAddressLineOne.String
	}
	if streetAddressLineTwo.Valid {
		vendorDetails.StreetAddressLineTwo = streetAddressLineTwo.String
	}
	if googleBusinessProfile.Valid {
		vendorDetails.GoogleBusinessProfile = googleBusinessProfile.String
	}

	return vendorDetails, nil
}

func GetSupplierDetails(supplierId string) (types.SupplierDetails, error) {
	query := `SELECT 
		s.supplier_id,
		s.name,
		s.membership_id,
		s.membership_cost::NUMERIC,
		s.membership_renewal,
		s.street_address_line_one,
		s.street_address_line_two,
		s.city_id,
		s.zip_code,
		s.state,
		s.google_business_profile
	FROM supplier AS s
	WHERE s.supplier_id = $1`

	var supplierDetails types.SupplierDetails

	row := DB.QueryRow(query, supplierId)

	var (
		membershipID          sql.NullString
		membershipCost        sql.NullFloat64
		membershipRenewal     sql.NullString
		streetAddressLineOne  sql.NullString
		streetAddressLineTwo  sql.NullString
		googleBusinessProfile sql.NullString
	)

	err := row.Scan(
		&supplierDetails.SupplierID,
		&supplierDetails.Name,
		&membershipID,
		&membershipCost,
		&membershipRenewal,
		&streetAddressLineOne,
		&streetAddressLineTwo,
		&supplierDetails.CityID,
		&supplierDetails.ZipCode,
		&supplierDetails.State,
		&googleBusinessProfile,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return supplierDetails, fmt.Errorf("no supplier found with ID %s", supplierId)
		}
		return supplierDetails, fmt.Errorf("error scanning row: %w", err)
	}

	// Handle nullable fields
	if membershipID.Valid {
		supplierDetails.MembershipID = membershipID.String
	}
	if membershipCost.Valid {
		supplierDetails.MembershipCost = membershipCost.Float64
	}
	if membershipRenewal.Valid {
		supplierDetails.MembershipRenewal = membershipRenewal.String
	}
	if streetAddressLineOne.Valid {
		supplierDetails.StreetAddressLineOne = streetAddressLineOne.String
	}
	if streetAddressLineTwo.Valid {
		supplierDetails.StreetAddressLineTwo = streetAddressLineTwo.String
	}
	if googleBusinessProfile.Valid {
		supplierDetails.GoogleBusinessProfile = googleBusinessProfile.String
	}

	return supplierDetails, nil
}

func GetBusinessDetails(businessID string) (types.BusinessDetails, error) {
	query := `SELECT 
		b.business_id,
		b.name,
		b.is_active,
		b.website,
		b.industry,
		b.google_business_profile
	FROM business AS b
	WHERE b.business_id = $1`

	var businessDetails types.BusinessDetails

	row := DB.QueryRow(query, businessID)

	err := row.Scan(
		&businessDetails.BusinessID,
		&businessDetails.Name,
		&businessDetails.IsActive,
		&businessDetails.Website,
		&businessDetails.Industry,
		&businessDetails.GoogleBusinessProfile,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return businessDetails, fmt.Errorf("no business found with ID %s", businessID)
		}
		return businessDetails, fmt.Errorf("error scanning row: %w", err)
	}

	return businessDetails, nil
}

func GetLocationsByBusiness(businessId string) ([]types.LocationList, error) {
	var locations []types.LocationList

	rows, err := DB.Query(`
		SELECT l.location_id, l.business_id, vl.location_type, l.name, l.longitude, l.latitude, l.street_address_line_one, l.street_address_line_two,
		c.name, l.zip_code, l.state, l.opening, l.closing 
		FROM "location" AS l
		JOIN "city" AS c ON c.city_id = l.city_id
		JOIN "vending_location" AS vl ON vl.vending_location_id = l.vending_location_id
		WHERE l.business_id = $1
	`, businessId)
	if err != nil {
		return locations, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var loc types.LocationList
		var streetAddressLineOne, streetAddressLineTwo, zipCode, state, opening, closing sql.NullString
		var longitude, latitude sql.NullFloat64

		err := rows.Scan(
			&loc.LocationID,
			&loc.BusinessID,
			&loc.VendingLocationType,
			&loc.Name,
			&longitude,
			&latitude,
			&streetAddressLineOne,
			&streetAddressLineTwo,
			&loc.City,
			&zipCode,
			&state,
			&opening,
			&closing,
		)
		if err != nil {
			return locations, fmt.Errorf("error scanning row: %w", err)
		}

		if longitude.Valid {
			loc.Longitude = longitude.Float64
		}
		if latitude.Valid {
			loc.Latitude = latitude.Float64
		}
		if streetAddressLineOne.Valid {
			loc.StreetAddressLineOne = streetAddressLineOne.String
		}
		if streetAddressLineTwo.Valid {
			loc.StreetAddressLineTwo = streetAddressLineTwo.String
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

func CreateMarketingImage(img models.Image) error {
	stmt, err := DB.Prepare(`
	INSERT INTO image (src, date_added, added_by_user_id)
	VALUES ($1, to_timestamp($2), $3)
`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(img.Src, img.DateAdded, img.AddedByUserID)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetMarketingImages() ([]string, error) {
	var images []string

	query := fmt.Sprintf(`SELECT '%s' || i.src AS url FROM "image" AS i;`, constants.AWSS3MarketingImagesPath)

	rows, err := DB.Query(query)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return images, err
	}
	defer rows.Close()

	for rows.Next() {

		var image string
		err := rows.Scan(
			&image,
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

func CreateMachine(form types.MachineForm) (int, error) {
	var machineId int
	stmt, err := DB.Prepare(`
		INSERT INTO machine (
			vending_type_id, 
			machine_status_id,
			vendor_id,
			year, 
			make, 
			model, 
			purchase_price, 
			purchase_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, to_timestamp($8))
		RETURNING machine_id
	`)
	if err != nil {
		return machineId, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	make := utils.CreateNullString(form.Make)
	model := utils.CreateNullString(form.Model)
	purchasePrice := utils.CreateNullFloat64(form.PurchasePrice)
	year := utils.CreateNullInt(form.Year)
	vendorID := utils.CreateNullInt(form.VendorID)
	purchaseDate := utils.CreateNullInt64(form.PurchaseDate)

	err = stmt.QueryRow(
		int64(*form.VendingTypeID),
		int64(*form.MachineStatusID),
		vendorID,
		year,
		make,
		model,
		purchasePrice,
		purchaseDate,
	).Scan(&machineId)

	if err != nil {
		return machineId, fmt.Errorf("error executing statement: %w", err)
	}

	return machineId, nil
}

func GetLocationDetails(businessID, locationID int) (types.LocationDetails, error) {
	query := `SELECT 
			l.location_id,
			l.business_id,
			l.vending_location_id,
			l.city_id,
			l.date_started,
			l.name,
			l.longitude,
			l.latitude,
			l.street_address_line_one,
			l.street_address_line_two,
			l.zip_code,
			l.state,
			l.opening,
			l.closing
		FROM location AS l
		WHERE l.location_id = $1 AND l.business_id = $2`

	var location types.LocationDetails

	row := DB.QueryRow(query, locationID, businessID)

	var dateStarted time.Time
	var streetAddressLineTwo, opening, closing sql.NullString
	var longitude, latitude sql.NullFloat64

	err := row.Scan(
		&location.LocationID,
		&location.BusinessID,
		&location.VendingLocationID,
		&location.CityID,
		&dateStarted,
		&location.Name,
		&longitude,
		&latitude,
		&location.StreetAddressLineOne,
		&streetAddressLineTwo,
		&location.ZipCode,
		&location.State,
		&opening,
		&closing,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return location, fmt.Errorf("no location found with ID %d", locationID)
		}
		return location, fmt.Errorf("error scanning row: %w", err)
	}

	location.DateStarted = dateStarted.Unix()

	if streetAddressLineTwo.Valid {
		location.StreetAddressLineTwo = streetAddressLineTwo.String
	}

	if opening.Valid {
		location.Opening = opening.String
	}

	if closing.Valid {
		location.Closing = closing.String
	}

	if longitude.Valid {
		location.Longitude = longitude.Float64
	}

	if latitude.Valid {
		location.Latitude = latitude.Float64
	}

	return location, nil
}

func GetMachinesByLocation(locationId int) ([]types.MachineList, error) {
	var machines []types.MachineList

	rows, err := DB.Query(`SELECT m.machine_id, CONCAT(m.year, ' ', m.make, ' ', m.model) AS machine_name,
	m.card_reader_serial_number, s.status, l.name, m.purchase_date
	FROM "machine" AS m
	JOIN machine_status AS s ON s.machine_status_id = m.machine_status_id
	JOIN location AS l ON l.location_id = m.location_id AND (l.location_id = $1 OR $1 IS NULL)
	ORDER BY m.purchase_date DESC;`, locationId)
	if err != nil {
		return machines, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var machine types.MachineList
		var dateCreated time.Time
		var location, cardReaderSerialNumber sql.NullString

		err := rows.Scan(
			&machine.MachineID,
			&machine.MachineName,
			&cardReaderSerialNumber,
			&machine.MachineStatus,
			&location,
			&dateCreated,
		)
		if err != nil {
			return machines, fmt.Errorf("error scanning row: %w", err)
		}

		if location.Valid {
			machine.Location = location.String
		}
		if cardReaderSerialNumber.Valid {
			machine.CardReaderSerialNumber = cardReaderSerialNumber.String
		}

		machine.PurchaseDate = utils.FormatDateMMDDYYYY(dateCreated.Unix())
		machines = append(machines, machine)
	}

	if err := rows.Err(); err != nil {
		return machines, fmt.Errorf("error iterating rows: %w", err)
	}

	return machines, nil
}

func GetMachineDetails(machineID int) (types.MachineDetails, error) {
	query := `SELECT
				m.machine_id,
				m.vending_type_id,
				m.machine_status_id,
				m.location_id,
				m.vendor_id,
				m.year,
				m.make,
				m.model,
				m.purchase_price::NUMERIC,
				m.purchase_date,
				m.card_reader_serial_number
			FROM machine AS m
			WHERE m.machine_id = $1`

	var machine types.MachineDetails
	var purchaseDate sql.NullTime
	var location sql.NullInt64
	var cardReaderSerialNumber sql.NullString

	row := DB.QueryRow(query, machineID)

	err := row.Scan(
		&machine.MachineID,
		&machine.VendingTypeID,
		&machine.MachineStatusID,
		&location,
		&machine.VendorID,
		&machine.Year,
		&machine.Make,
		&machine.Model,
		&machine.PurchasePrice,
		&purchaseDate,
		&cardReaderSerialNumber,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return machine, fmt.Errorf("no machine found with ID %d", machineID)
		}
		return machine, fmt.Errorf("error scanning row: %w", err)
	}

	// Check nullable fields
	if purchaseDate.Valid {
		machine.PurchaseDate = purchaseDate.Time.Unix()
	}
	if cardReaderSerialNumber.Valid {
		machine.CardReaderSerialNumber = cardReaderSerialNumber.String
	}
	if location.Valid {
		machine.LocationID = int(location.Int64)
	}

	return machine, nil
}

func CreateSeedLiveTransaction(transaction types.SeedLiveTransaction) error {
	stmt, err := DB.Prepare(`
		INSERT INTO seed_live_transaction (
			terminal_number, 
			transaction_ref_number, 
			transaction_type, 
			card_number, 
			total_amount, 
			vended_columns, 
			price, 
			mdb_number, 
			number_of_products_vended, 
			timestamp, 
			card_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, to_timestamp($10), $11)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var (
		terminalNumber         = transaction.TerminalNumber
		transactionRefNumber   = transaction.TransactionRefNumber
		transactionType        = transaction.TransactionType
		cardNumber             = utils.CreateNullString(&transaction.CardNumber)
		totalAmount            = transaction.TotalAmount
		vendedColumns          = transaction.VendedColumns
		price                  = transaction.Price
		mdbNumber              = utils.CreateNullInt(&transaction.MDBNumber)
		numberOfProductsVended = utils.CreateNullInt(&transaction.NumberOfProductsVended)
		timestamp              = transaction.Timestamp.Unix()
		cardId                 = utils.CreateNullString(&transaction.CardId)
	)

	_, err = stmt.Exec(
		terminalNumber,
		transactionRefNumber,
		transactionType,
		cardNumber,
		totalAmount,
		vendedColumns,
		price,
		mdbNumber,
		numberOfProductsVended,
		timestamp,
		cardId,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func GetProductCategories() ([]models.ProductCategory, error) {
	var productCategories []models.ProductCategory

	rows, err := DB.Query(`SELECT "product_category_id", "name" FROM "product_category"`)
	if err != nil {
		return productCategories, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var productCategory models.ProductCategory
		err := rows.Scan(&productCategory.ProductCategoryID, &productCategory.Name)
		if err != nil {
			return productCategories, fmt.Errorf("error scanning row: %w", err)
		}
		productCategories = append(productCategories, productCategory)
	}

	if err := rows.Err(); err != nil {
		return productCategories, fmt.Errorf("error iterating rows: %w", err)
	}

	return productCategories, nil
}

func GetProductList(pageNum int) ([]types.ProductList, int, error) {
	var products []types.ProductList
	var totalRows int

	var offset = (pageNum - 1) * int(constants.LeadsPerPage)

	rows, err := DB.Query(`
		SELECT 
			p.product_id,
			p.name,
			c.name AS category,
			p.size,
			p.size_type,
			p.upc,
			COUNT(*) OVER() AS total_rows
		FROM "product" AS p
		JOIN product_category AS c ON c.product_category_id = p.product_category_id
		ORDER BY p.product_id DESC
		LIMIT $1
		OFFSET $2;
	`, constants.LeadsPerPage, offset)
	if err != nil {
		return products, totalRows, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product types.ProductList

		var size sql.NullFloat64
		var sizeType, upc sql.NullString

		err := rows.Scan(
			&product.ProductID,
			&product.Name,
			&product.Category,
			&size,
			&sizeType,
			&upc,
			&totalRows,
		)
		if err != nil {
			return products, totalRows, fmt.Errorf("error scanning row: %w", err)
		}

		if size.Valid {
			product.Size = size.Float64
		}

		if sizeType.Valid {
			product.SizeType = sizeType.String
		}

		if upc.Valid {
			product.UPC = upc.String
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return products, totalRows, fmt.Errorf("error iterating rows: %w", err)
	}

	return products, totalRows, nil
}

func CreateProduct(form types.ProductForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO product (
			name,
			size,
			size_type,
			upc,
			product_category_id
		) VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	name := utils.CreateNullString(form.Name)
	productCategoryID := utils.CreateNullInt(form.ProductCategoryID)
	size := utils.CreateNullFloat64(form.Size)
	sizeType := utils.CreateNullString(form.SizeType)
	upc := utils.CreateNullString(form.UPC)

	_, err = stmt.Exec(
		name,
		size,
		sizeType,
		upc,
		productCategoryID,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateProduct(productId int, form types.ProductForm) error {
	stmt, err := DB.Prepare(`
		UPDATE product 
		SET 
		name = $2,
			size = $3,
			size_type = $4,
			upc = $5,
			product_category_id = $6
		WHERE product_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	name := utils.CreateNullString(form.Name)
	productCategoryID := utils.CreateNullInt(form.ProductCategoryID)
	size := utils.CreateNullFloat64(form.Size)
	sizeType := utils.CreateNullString(form.SizeType)
	upc := utils.CreateNullString(form.UPC)

	_, err = stmt.Exec(
		productId,
		name,
		size,
		sizeType,
		upc,
		productCategoryID,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func DeleteProduct(id int) error {
	sqlStatement := `
        DELETE FROM product WHERE product_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func GetProductDetails(productID string) (types.ProductDetails, error) {
	query := `SELECT 
		p.product_id,
		p.name,
		p.size,
		p.size_type,
		p.product_category_id,
		p.upc
	FROM product AS p
	WHERE p.product_id = $1`

	var productDetails types.ProductDetails

	row := DB.QueryRow(query, productID)

	var (
		size     sql.NullFloat64
		sizeType sql.NullString
		upc      sql.NullString
	)

	err := row.Scan(
		&productDetails.ProductID,
		&productDetails.Name,
		&size,
		&sizeType,
		&productDetails.ProductCategoryID,
		&upc,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return productDetails, fmt.Errorf("no product found with ID %s", productID)
		}
		return productDetails, fmt.Errorf("error scanning row: %w", err)
	}

	if upc.Valid {
		productDetails.UPC = upc.String
	}

	if size.Valid {
		productDetails.Size = size.Float64
	}

	if sizeType.Valid {
		productDetails.SizeType = sizeType.String
	}

	return productDetails, nil
}

func GetProductBatchList(productId string) ([]types.ProductBatchList, error) {
	var batches []types.ProductBatchList

	rows, err := DB.Query(`
		SELECT 
			pb.product_batch_id,
			pb.product_id,
			s.name,
			pb.expiration_date,
			pb.unit_cost::NUMERIC,
			pb.quantity,
			pb.date_purchased
		FROM "product_batch" AS pb
		JOIN supplier AS s ON pb.supplier_id = s.supplier_id
		WHERE pb.product_id = $1
		ORDER BY pb.date_purchased ASC;
	`, productId)
	if err != nil {
		return batches, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var batch types.ProductBatchList
		var datePurchased, expirationDate time.Time

		err := rows.Scan(
			&batch.ProductBatchID,
			&batch.ProductID,
			&batch.Supplier,
			&expirationDate,
			&batch.UnitCost,
			&batch.Quantity,
			&datePurchased,
		)
		if err != nil {
			return batches, fmt.Errorf("error scanning row: %w", err)
		}

		batch.DatePurchased = utils.FormatDateMMDDYYYY(datePurchased.Unix())
		batch.ExpirationDate = utils.FormatDateMMDDYYYY(expirationDate.Unix())

		batches = append(batches, batch)
	}

	if err := rows.Err(); err != nil {
		return batches, fmt.Errorf("error iterating rows: %w", err)
	}

	return batches, nil
}

func CreateProductBatch(form types.ProductBatchForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO product_batch (
			product_id,
			supplier_id,
			expiration_date,
			date_purchased,
			unit_cost,
			quantity
		) VALUES ($1, $2, to_timestamp($3), to_timestamp($4), $5, $6)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	productId := utils.CreateNullInt(form.ProductID)
	supplierId := utils.CreateNullInt(form.SupplierID)
	unitCost := utils.CreateNullFloat64(form.UnitCost)
	quantity := utils.CreateNullInt(form.Quantity)

	_, err = stmt.Exec(
		productId,
		supplierId,
		form.ExpirationDate,
		form.DatePurchased,
		unitCost,
		quantity,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateProductBatch(productId int, form types.ProductBatchForm) error {
	stmt, err := DB.Prepare(`
		UPDATE product_batch 
		SET 
			supplier_id = $2,
			expiration_date = to_timestamp($3),
			date_purchased = to_timestamp($4),
			unit_cost = $5,
			quantity = $6
		WHERE product_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	supplierId := utils.CreateNullInt(form.SupplierID)
	unitCost := utils.CreateNullFloat64(form.UnitCost)
	quantity := utils.CreateNullInt(form.Quantity)

	_, err = stmt.Exec(
		productId,
		supplierId,
		form.ExpirationDate,
		form.DatePurchased,
		unitCost,
		quantity,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func DeleteProductBatch(id int) error {
	sqlStatement := `
        DELETE FROM product_batch WHERE product_batch_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func GetProductBatchDetails(productId, productBatchId int) (models.ProductBatch, error) {
	var productBatch models.ProductBatch

	row := DB.QueryRow(`
		SELECT 
			pb.product_batch_id,
			pb.product_id,
			pb.expiration_date,
			pb.unit_cost::NUMERIC,
			pb.quantity,
			pb.date_purchased,
			pb.supplier_id
		FROM "product_batch" AS pb
		WHERE pb.product_id = $1 AND pb.product_batch_id = $2;
	`, productId, productBatchId)

	var datePurchased, expirationDate time.Time

	err := row.Scan(
		&productBatch.ProductBatchID,
		&productBatch.ProductID,
		&expirationDate,
		&productBatch.UnitCost,
		&productBatch.Quantity,
		&datePurchased,
		&productBatch.SupplierID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return productBatch, fmt.Errorf("no product batch found with ID %d", productBatchId)
		}
		return productBatch, fmt.Errorf("error scanning row: %w", err)
	}

	productBatch.DatePurchased = datePurchased.Unix()
	productBatch.ExpirationDate = expirationDate.Unix()

	return productBatch, nil
}

func GetMachineSlotsByMachineID(machineId string) ([]types.SlotList, error) {
	var slots []types.SlotList

	rows, err := DB.Query(`
	SELECT 
		s.slot_id,
		s.slot,
		s.machine_id,
		s.machine_code,
		s.price::NUMERIC,
		s.capacity,
		MAX(r.date_refilled),
		r.refill_id
	FROM "slot" AS s
	LEFT JOIN refill AS r ON r.slot_id = s.slot_id
	WHERE s.machine_id = $1
	GROUP BY s.slot_id, s.slot, s.machine_id, s.machine_code, s.price, s.capacity, r.refill_id
	ORDER BY s.slot ASC;
	`, machineId)
	if err != nil {
		return slots, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var slot types.SlotList

		var dateRefilled sql.NullTime

		err := rows.Scan(
			&slot.SlotID,
			&slot.Slot,
			&slot.MachineID,
			&slot.MachineCode,
			&slot.Price,
			&slot.Capacity,
			&dateRefilled,
			&slot.LastRefillID,
		)
		if err != nil {
			return slots, fmt.Errorf("error scanning row: %w", err)
		}

		if dateRefilled.Valid {
			slot.LastRefill = utils.FormatTimestamp(dateRefilled.Time.Unix())
		}

		slots = append(slots, slot)
	}

	if err := rows.Err(); err != nil {
		return slots, fmt.Errorf("error iterating rows: %w", err)
	}

	return slots, nil
}

func GetAvailableProductBatches() ([]types.AvailableProductBatches, error) {
	var batches []types.AvailableProductBatches

	rows, err := DB.Query(`
		SELECT 
			pb.product_batch_id,
			s.name,
			p.name,
			pb.unit_cost::NUMERIC,
			pb.quantity,
			pb.date_purchased,
			pb.expiration_date
		FROM "product_batch" AS pb
		JOIN supplier AS s ON pb.supplier_id = s.supplier_id
		JOIN product AS p ON pb.product_id = p.product_id
		WHERE pb.quantity > 0 AND pb.expiration_date > CURRENT_DATE
		ORDER BY pb.date_purchased ASC;
	`)
	if err != nil {
		return batches, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var batch types.AvailableProductBatches
		var datePurchased, expirationDate time.Time

		err := rows.Scan(
			&batch.ProductBatchID,
			&batch.Supplier,
			&batch.ProductName,
			&batch.UnitCost,
			&batch.Quantity,
			&datePurchased,
			&expirationDate,
		)
		if err != nil {
			return batches, fmt.Errorf("error scanning row: %w", err)
		}

		batch.DatePurchased = utils.FormatDateMMDDYYYY(datePurchased.Unix())
		batch.ExpirationDate = utils.FormatDateMMDDYYYY(expirationDate.Unix())

		batches = append(batches, batch)
	}

	if err := rows.Err(); err != nil {
		return batches, fmt.Errorf("error iterating rows: %w", err)
	}

	return batches, nil
}

func CreateSlot(form types.SlotForm) (int, error) {
	var slotId int
	stmt, err := DB.Prepare(`
		INSERT INTO slot (
			nickname,
			slot,
			machine_code,
			machine_id,
			price,
			capacity
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING slot_id
	`)
	if err != nil {
		return slotId, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	nickname := utils.CreateNullString(form.Nickname)
	slot := utils.CreateNullString(form.Slot)
	machineCode := utils.CreateNullString(form.MachineCode)
	machineID := utils.CreateNullInt(form.MachineID)
	price := utils.CreateNullFloat64(form.Price)
	capacity := utils.CreateNullInt(form.Capacity)

	err = stmt.QueryRow(
		nickname,
		slot,
		machineCode,
		machineID,
		price,
		capacity,
	).Scan(&slotId)

	if err != nil {
		return slotId, fmt.Errorf("error executing statement: %w", err)
	}

	return slotId, nil
}

func UpdateSlot(slotId int, form types.SlotForm) error {
	stmt, err := DB.Prepare(`
		UPDATE slot 
		SET 
			nickname = COALESCE($2, nickname),
			slot = COALESCE($3, slot),
			machine_code = COALESCE($4, machine_code),
			machine_id = COALESCE($5, machine_id),
			price = COALESCE($6, price),
			capacity = COALESCE($7, capacity)
		WHERE slot_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	nickname := utils.CreateNullString(form.Nickname)
	slot := utils.CreateNullString(form.Slot)
	machineCode := utils.CreateNullString(form.MachineCode)
	machineID := utils.CreateNullInt(form.MachineID)
	price := utils.CreateNullFloat64(form.Price)
	capacity := utils.CreateNullInt(form.Capacity)

	_, err = stmt.Exec(
		slotId,
		nickname,
		slot,
		machineCode,
		machineID,
		price,
		capacity,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func DeleteSlot(id int) error {
	sqlStatement := `
        DELETE FROM slot WHERE slot_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func GetSlotDetails(machineId, slotId string) (types.SlotDetails, error) {
	query := `SELECT 
		s.slot_id,
		s.nickname,
		s.slot,
		s.machine_code,
		s.machine_id,
		s.price,
		s.capacity
	FROM slot AS s
	WHERE s.slot_id = $1 AND s.machine_id = $2`

	var businessDetails types.SlotDetails

	row := DB.QueryRow(query, slotId, machineId)

	err := row.Scan(
		&businessDetails.SlotID,
		&businessDetails.Nickname,
		&businessDetails.Slot,
		&businessDetails.MachineCode,
		&businessDetails.MachineID,
		&businessDetails.Price,
		&businessDetails.Capacity,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return businessDetails, fmt.Errorf("no business found with ID %s", slotId)
		}
		return businessDetails, fmt.Errorf("error scanning row: %w", err)
	}

	return businessDetails, nil
}

func GetProductSlotAssignments(slotId string) ([]types.ProductSlotAssignment, error) {
	query := `SELECT 
		psa.product_slot_assignment_id,
		s.slot,
		p.name,
		psa.date_assigned
	FROM product_slot_assignment AS psa
	JOIN slot AS s ON psa.slot_id = s.slot_id
	JOIN product_batch AS pb ON pb.product_batch_id = psa.product_batch_id
	JOIN product AS p ON p.product_id = pb.product_id
	WHERE psa.slot_id = $1`

	var productSlotAssignments []types.ProductSlotAssignment

	rows, err := DB.Query(query, slotId)

	if err != nil {
		return productSlotAssignments, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var assignment types.ProductSlotAssignment
		var dateAssigned time.Time

		err := rows.Scan(
			&assignment.ProductSlotAssignmentID,
			&assignment.Slot,
			&assignment.Product,
			&dateAssigned,
		)
		if err != nil {
			return productSlotAssignments, fmt.Errorf("error scanning row: %w", err)
		}

		assignment.DateAssigned = utils.FormatDateMMDDYYYY(dateAssigned.Unix())

		productSlotAssignments = append(productSlotAssignments, assignment)
	}

	if err := rows.Err(); err != nil {
		return productSlotAssignments, fmt.Errorf("error iterating rows: %w", err)
	}

	return productSlotAssignments, nil
}

func CreateProductSlotAssignment(form types.ProductSlotAssignmentForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO product_slot_assignment (
			slot_id,
			product_batch_id,
			date_assigned
		) VALUES ($1, $2, to_timestamp($3))
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		form.SlotID,
		form.ProductBatchID,
		form.DateAssigned,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func DeleteProductSlotAssignment(id int) error {
	sqlStatement := `
        DELETE FROM product_slot_assignment WHERE product_slot_assignment_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProductSlotAssignment(form types.ProductSlotAssignmentForm) error {
	stmt, err := DB.Prepare(`
		UPDATE product_slot_assignment 
		SET 
			slot_id = COALESCE($2, slot_id),
			product_batch_id = COALESCE($3, product_batch_id),
			date_assigned = COALESCE(to_timestamp($4), date_assigned)
		WHERE product_slot_assignment_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	slotID := utils.CreateNullInt(form.SlotID)
	productBatchID := utils.CreateNullInt(form.ProductBatchID)
	dateAssigned := utils.CreateNullInt64(form.DateAssigned)

	_, err = stmt.Exec(
		slotID,
		productBatchID,
		dateAssigned,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CreateRefill(form types.RefillForm) error {
	stmt, err := DB.Prepare(`
		INSERT INTO refill (
			slot_id,
			date_refilled
		) VALUES ($1, to_timestamp($2))
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		form.SlotID,
		form.DateRefilled,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func DeleteRefill(id int) error {
	sqlStatement := `
        DELETE FROM refill WHERE refill_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func CreateMachineLocationAssignment(form models.MachineLocationAssignment) error {
	stmt, err := DB.Prepare(`
		INSERT INTO machine_location_assignment (
			location_id,
			machine_id,
			date_assigned
		) VALUES ($1, $2, to_timestamp($3))
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		form.LocationID,
		form.MachineID,
		form.DateAssigned,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CreateMachineCardReaderAssignment(form models.MachineCardReaderAssignment) error {
	stmt, err := DB.Prepare(`
		INSERT INTO machine_card_reader (
			card_reader_serial_number,
			machine_id,
			date_assigned
		) VALUES ($1, $2, to_timestamp($3))
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		form.CardReaderSerialNumber,
		form.MachineID,
		form.DateAssigned,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CreateSlotPriceLog(form models.SlotPriceLog) error {
	stmt, err := DB.Prepare(`
		INSERT INTO slot_price_log (
			slot_id,
			date_assigned
		) VALUES ($1, $2, to_timestamp($3))
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		form.SlotID,
		form.DateAssigned,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func DeleteMachineLocationAssignment(id int) error {
	sqlStatement := `
        DELETE FROM machine_location_assignment WHERE machine_location_assignment_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteMachineCardReaderAssignment(id int) error {
	sqlStatement := `
        DELETE FROM machine_card_reader WHERE machine_card_reader_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSlotPriceLog(id int) error {
	sqlStatement := `
        DELETE FROM slot_price_log WHERE slot_price_log_id = $1
    `
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateMachineLocationAssignment(form models.MachineLocationAssignment) error {
	stmt, err := DB.Prepare(`
		UPDATE machine_location_assignment
		SET location_id = COALESCE($2, location_id),
			machine_id = COALESCE($3, machine_id),
			date_assigned = COALESCE(to_timestamp($4), date_assigned)
		WHERE machine_location_assignment_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		form.MachineLocationAssignmentID,
		form.LocationID,
		form.MachineID,
		form.DateAssigned,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateMachineCardReaderAssignment(form models.MachineCardReaderAssignment) error {
	stmt, err := DB.Prepare(`
		UPDATE machine_card_reader
		SET card_reader_serial_number = COALESCE($2, card_reader_serial_number),
			machine_id = COALESCE($3, machine_id),
			date_assigned = COALESCE(to_timestamp($4), date_assigned)
		WHERE machine_card_reader_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		form.MachineCardReaderID,
		form.CardReaderSerialNumber,
		form.MachineID,
		form.DateAssigned,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func UpdateSlotPriceLog(form models.SlotPriceLog) error {
	stmt, err := DB.Prepare(`
		UPDATE slot_price_log
		SET slot_id = COALESCE($2, slot_id),
			date_assigned = COALESCE($3, date_assigned)
		WHERE slot_price_log_id = $1
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		form.SlotPriceLogID,
		form.SlotID,
		form.DateAssigned,
	)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}
