package github

import (
	"context"
	"os"

	g "github.com/google/go-github/v60/github"
)

type GitHubClientConfig struct {
	context context.Context
}

type GitHubClient struct {
	*GitHubClientConfig
	client *g.Client
}

func GithubAuthToken() string {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token
	}
	return ""
}

func NewGitHubClient(config *GitHubClientConfig) (*GitHubClient, error) {
	if config.context == nil {
		config.context = context.Background()
	}

	client := g.NewClient(nil).WithAuthToken(GithubAuthToken())
	return &GitHubClient{
		config,
		client,
	}, nil
}

func (gc *GitHubClient) CreateIssue(owner, repo, title, body string) error {
	_, _, err := gc.client.Repositories.Get(gc.context, owner, repo)
	if err != nil {
		return err
	}

	_, _, err = gc.client.Issues.Create(gc.context, owner, repo,
		&g.IssueRequest{
			Title: &title,
			Body:  &body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
