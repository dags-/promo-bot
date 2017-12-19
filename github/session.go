package github

import (
	"golang.org/x/oauth2"
	"context"
	"github.com/google/go-github/github"
)

type Session struct {
	Client *github.Client
	Ctx    context.Context
}

type Repo struct {
	session *Session
	owner   string
	name    string
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
		Ctx: ctx,
	}
}

func (s *Session) NewRepo(owner, name string) (Repo) {
	return Repo{
		session: s,
		owner: owner,
		name: name,
		ref: "master",
	}
}

func (s *Session) do(method, path string, requestBody, responseBody interface{}) (error)  {
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