package services

import (
	"context"
	"fmt"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

func GetLatestGithubCommit(repo, branch string) (*github.RepositoryCommit, error) {
	ctx := context.Background()

	var client *github.Client

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: constants.GithubAccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	commits, _, err := client.Repositories.ListCommits(ctx, constants.GithubOwnerName, repo, &github.CommitsListOptions{
		SHA: branch,
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching commits: %w", err)
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits found in the repository")
	}

	return commits[0], nil
}
