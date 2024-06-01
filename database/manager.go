package database

func InsertCSRFToken(token CSRFToken) error {
	stmt, err := db.Prepare("INSERT INTO csrf_token(expiry_time, token) VALUES(?, ?)")
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

func GetCSRFToken(decryptedToken string) (CSRFToken, error) {
	stmt, err := db.Prepare("SELECT * FROM csrf_token WHERE token = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(decryptedToken)

	var token CSRFToken
	err = row.Scan(&token.CSRFTokenID, &token.ExpiryTime, &token.Token, &token.IsUsed)
	if err != nil {
		return nil, err
	}

	return token, nil
}
