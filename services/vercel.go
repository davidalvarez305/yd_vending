package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/davidalvarez305/yd_vending/types"
)

func CreateVercelProject(slug, teamID, token string, project types.VercelProjectRequestBody) error {
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
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("VERCEL CREATION ERROR: %s\n", bodyString)
		return fmt.Errorf("error creating project, received status code: %d", resp.StatusCode)
	}

	log.Printf("Successfully created Vercel project with status code %d", resp.StatusCode)

	return nil
}

func UpdateVercelProject(slug, teamID, token string, project types.VercelProjectRequestBody) error {
	url := fmt.Sprintf("https://api.vercel.com/v10/projects?slug=%s&teamId=%s", slug, teamID)

	// Create the JSON body from the struct
	body, err := json.Marshal(project)
	if err != nil {
		return fmt.Errorf("error marshalling project data: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
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
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("VERCEL UPDATING ERROR: %s\n", bodyString)
		return fmt.Errorf("error updating project, received status code: %d", resp.StatusCode)
	}

	log.Printf("Successfully created Vercel project with status code %d", resp.StatusCode)

	return nil
}

func DeleteVercelProject(slug, teamID, token, projectId string) error {
	url := fmt.Sprintf("https://api.vercel.com/v10/projects/%s?slug=%s&teamId=%s", projectId, slug, teamID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("VERCEL DELETING ERROR: %s\n", bodyString)
		return fmt.Errorf("error deleting project, received status code: %d", resp.StatusCode)
	}

	log.Printf("Successfully created Vercel project with status code %d", resp.StatusCode)

	return nil
}
