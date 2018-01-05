package github

import (
	"encoding/base64"
	"fmt"
)

type Contents struct {
	Name    string `json:"name"`
	Sha     string `json:"sha"`
	URL     string `json:"download_url"`
	Message string `json:"message"`
}

type FileCreate struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Content string `json:"content"`
	Branch  string `json:"branch"`
}

type FileUpdate struct {
	FileCreate
	Sha string `json:"sha"`
}

func (r *Repo) GetContents() ([]Contents, error) {
	var contents []Contents

	session := r.session
	url := fmt.Sprintf("repos/%s/%s/contents", r.owner, r.name)
	req, err := session.Client.NewRequest("GET", url, nil)
	if err != nil {
		return contents, err
	}

	_, err = session.Client.Do(session.Ctx, req, &contents)
	return contents, err
}

func (b *Branch) CreateFile(file string, content []byte) (error) {
	r := b.Repo
	s := r.session

	encoded := base64.StdEncoding.EncodeToString(content)
	path := fmt.Sprintf("repos/%s/%s/contents/%s", r.owner, r.name, file)

	var body interface{}
	create := FileCreate{
		Path: path,
		Message: "This is a test",
		Content: encoded,
		Branch: b.Name,
	}

	body = create

	if query, err := b.getFile(file); err == nil {
		fmt.Println("File exists ", file)
		var update FileUpdate
		update.FileCreate = create
		update.Sha = query.Sha
		body = update // file exists so change body to FileUpdate type
	}

	return s.do("PUT", path, &body, nil)
}

func (b *Branch) getFile(file string) (Contents, error) {
	var contents Contents
	r := b.Repo
	s := r.session
	path := fmt.Sprintf("repos/%s/%s/contents/%s?ref=%s", r.owner, r.name, file, b.Name)
	err := s.do("GET", path, nil, &contents)
	return contents, err
}