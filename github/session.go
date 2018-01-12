package github

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Session struct {
	Client *github.Client
	Ctx    context.Context
}

type Repo struct {
	session *Session
	Owner   string
	Name    string
	ref     string
}

func NewSession(token string) Session {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	auth := oauth2.NewClient(ctx, ts)
	client := github.NewClient(auth)
	return Session{
		Client: client,
		Ctx:    ctx,
	}
}

func (s *Session) NewRepo(owner, name string) (Repo) {
	return Repo{
		session: s,
		Owner:   owner,
		Name:    name,
		ref:     "master",
	}
}

func (s *Session) RemainingRate() (int, error) {
	rate, _, err := s.Client.RateLimits(s.Ctx)
	if err != nil {
		return 0, err
	}
	return rate.Core.Remaining, nil
}

func (s *Session) do(method, path string, requestBody, responseBody interface{}) (error) {
	request, err := s.Client.NewRequest(method, path, requestBody)
	if err != nil {
		return err
	}

	if request.Body != nil {
		defer request.Body.Close()
	}

	response, err := s.Client.Do(s.Ctx, request, responseBody)
	if err != nil {
		return err
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	return nil // great success!
}
