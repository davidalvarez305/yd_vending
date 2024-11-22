package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/davidalvarez305/yd_vending/types"
)

func CreateVercelProject(slug, teamID, token string, project types.VercelProjectRequestBody) (types.CreateVercelProjectResponse, error) {
	url := fmt.Sprintf("https://api.vercel.com/v10/projects?slug=%s&teamId=%s", slug, teamID)
	var createProjectResp types.CreateVercelProjectResponse

	body, err := json.Marshal(project)
	if err != nil {
		return createProjectResp, fmt.Errorf("error marshalling project data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return createProjectResp, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return createProjectResp, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return createProjectResp, fmt.Errorf("failed to read response body: %w", err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("VERCEL CREATION ERROR: %s\n", bodyString)
		return createProjectResp, fmt.Errorf("error creating project, received status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&createProjectResp); err != nil {
		return createProjectResp, fmt.Errorf("error decoding response: %w", err)
	}

	log.Printf("Successfully created Vercel project with status code %d", resp.StatusCode)

	return createProjectResp, nil
}

func UpdateVercelProject(slug, teamID, token string, project types.VercelProjectRequestBody) (types.UpdateVercelProjectResponse, error) {
	url := fmt.Sprintf("https://api.vercel.com/v10/projects?slug=%s&teamId=%s", slug, teamID)
	var updateProjectResp types.UpdateVercelProjectResponse

	body, err := json.Marshal(project)
	if err != nil {
		return updateProjectResp, fmt.Errorf("error marshalling project data: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return updateProjectResp, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return updateProjectResp, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return updateProjectResp, fmt.Errorf("failed to read response body: %w", err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("VERCEL UPDATING ERROR: %s\n", bodyString)
		return updateProjectResp, fmt.Errorf("error updating project, received status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&updateProjectResp); err != nil {
		return updateProjectResp, fmt.Errorf("error decoding response: %w", err)
	}

	log.Printf("Successfully updated Vercel project with status code %d", resp.StatusCode)

	return updateProjectResp, nil
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

func GetVercelEnvironmentVariables(token, projectId, gitBranch, slug, teamId string) (types.GetVercelEnvironmentVariablesResponse, error) {
	url := fmt.Sprintf(
		"https://api.vercel.com/v9/projects/%s/env?decrypt=true&gitBranch=%s&slug=%s&source=vercel-cli:pull&teamId=%s",
		projectId, gitBranch, slug, teamId,
	)

	var response types.GetVercelEnvironmentVariablesResponse

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return response, fmt.Errorf("failed to read response body: %w", err)
		}
		bodyString := string(bodyBytes)
		log.Printf("Error from Vercel: %s", bodyString)
		return response, fmt.Errorf("error fetching environment variables, received status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("error decoding response: %w", err)
	}

	log.Printf("Successfully fetched Vercel environment variables with status code %d", resp.StatusCode)

	return response, nil
}

func UpdateVercelEnvironmentVariables(token, projectId, gitBranch, slug, teamId string, body []types.UpdateVercelEnvironmentVariablesBody) error {
	var wg sync.WaitGroup

	for _, variable := range body {
		wg.Add(1)

		go func(variable types.UpdateVercelEnvironmentVariablesBody) {
			defer wg.Done()

			url := fmt.Sprintf(
				"https://api.vercel.com/v9/projects/%s/env/%s?decrypt=true&gitBranch=%s&slug=%s&source=vercel-cli:pull&teamId=%s",
				variable.ID, projectId, gitBranch, slug, teamId,
			)

			req, err := http.NewRequest("PATCH", url, nil)
			if err != nil {
				log.Printf("Error creating request for variable %s: %v", variable.ID, err)
				return
			}

			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{
				Timeout: 30 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error sending request for variable %s: %v", variable.ID, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Failed to read response body for variable %s: %v", variable.ID, err)
					return
				}
				bodyString := string(bodyBytes)
				log.Printf("Error from Vercel for variable %s: %s", variable.ID, bodyString)
				return
			}
		}(variable)
	}

	wg.Wait()

	return nil
}

func CreateVercelProjectDeployment(slug, teamID, token string, deploymentBody types.DeployVercelProjectBody) error {
	url := fmt.Sprintf("https://api.vercel.com/v13/deployments?forceNew=1&skipAutoDetectionConfirmation=1&slug=%s&teamId=%s", slug, teamID)

	body, err := json.Marshal(deploymentBody)
	if err != nil {
		return fmt.Errorf("error marshalling project data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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
		fmt.Printf("VERCEL DEPLOYMENT ERROR: %s\n", bodyString)
		return fmt.Errorf("error deploying project, received status code: %d", resp.StatusCode)
	}

	return nil
}
