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

func ManageVercelProject(method, slug, teamID, token string, project types.VercelProjectRequestBody, projectId string) error {
	var url string
	if method == "POST" {
		url = fmt.Sprintf("https://api.vercel.com/v10/projects?slug=%s&teamId=%s", slug, teamID)
	} else {
		url = fmt.Sprintf("https://api.vercel.com/v10/projects/%s?slug=%s&teamId=%s", projectId, slug, teamID)
	}

	var body []byte
	if method != "DELETE" {
		var err error
		body, err = json.Marshal(project)
		if err != nil {
			return fmt.Errorf("error marshalling project data: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
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
		fmt.Printf("VERCEL ERROR: %s\n", bodyString)
		return fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	log.Printf("Successfully completed Vercel %s request with status code %d", method, resp.StatusCode)
	return nil
}
