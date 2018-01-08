package github

import (
	"fmt"
)

type PRCreate struct {
	Title  string `json:"title"`
	Head   string `json:"head"`
	Base   string `json:"base"`
	Body   string `json:"body"`
	Modify bool   `json:"maintainer_can_modify"`
}

type PRResponse struct {
	URL     string `json:"html_url"`
	Message string `json:"message"`
}

func (b *Branch) CreatePR(title, comment string) (PRResponse, error) {
	var response PRResponse

	r := b.Repo
	s := r.session
	path := fmt.Sprintf("repos/%s/%s/pulls", r.owner, r.name)
	body := PRCreate{
		Title:  title,
		Head:   b.Name,
		Base:   r.ref,
		Body:   comment,
		Modify: false,
	}

	return response, s.do("POST", path, &body, &response)
}
