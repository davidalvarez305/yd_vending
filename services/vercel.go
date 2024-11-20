package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davidalvarez305/yd_vending/types"
)

func CreateVercelProject(slug, teamID, token string, project types.CreateVercelProjectBody) error {
	url := fmt.Sprintf("https://api.vercel.com/v10/projects?slug=%s&teamId=%s", slug, teamID)

	// Create the JSON body from the struct
	body, err := json.Marshal(project)
	if err != nil {
		return fmt.Errorf("error marshalling project data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Add the Authorization header
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("error creating project, received status code: %d", resp.StatusCode)
	}

	log.Printf("Successfully created Vercel project with status code %d", resp.StatusCode)

	return nil
}
