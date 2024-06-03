package database

import (
	"fmt"

	"github.com/davidalvarez305/budgeting/models"
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
