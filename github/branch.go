package github

import (
	"fmt"
)

type Branch struct {
	Repo *Repo
	Name string
}

type BranchQuery struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type BranchCreate struct {
	Ref string `json:"ref"`
	Sha string `json:"sha"`
}

func (r *Repo) CreateBranch(name string) (Branch, error) {
	var branch Branch
	s := r.session

	if _, err := r.branchExists(name); err == nil {
		branch.Repo = r
		branch.Name = name
		return branch, nil
	}

	sha, _, err := s.Client.Repositories.GetCommitSHA1(s.Ctx, r.Owner, r.Name, r.ref, "")
	if err != nil {
		return branch, err
	}

	path := fmt.Sprintf("repos/%s/%s/git/refs", r.Owner, r.Name)
	body := BranchCreate{Ref: fmt.Sprintf("refs/heads/%s", name), Sha: sha}
	err = s.do("POST", path, &body, &branch)

	if err == nil {
		branch.Repo = r
		branch.Name = name
	}

	return branch, err
}

func (r *Repo) branchExists(name string) (BranchQuery, error) {
	var branch BranchQuery
	path := fmt.Sprintf("repos/%s/%s/branches/%s", r.Owner, r.Name, name)
	err := r.session.do("GET", path, nil, &branch)
	if branch.Name == name {
		return branch, nil
	}
	return branch, err
}
