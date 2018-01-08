package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dags-/promo-bot/github"
	"github.com/dags-/promo-bot/promo"
	"github.com/dags-/promo-bot/util"
	"github.com/qiangxue/fasthttp-routing"
)

func (s *Server) handleGet(c *routing.Context) error {
	id := c.Param("auth")
	if s.auth.isRateLimited(id) {
		return errors.New("please come back later")
	}

	if !s.auth.isAuthenticated(id) {
		return s.redirect(c)
	}

	c.Response.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Response.Header.Set("Pragma", "no-cache")
	c.Response.Header.Set("Expires", "0")
	c.SendFile("_public/form/index.html")

	return nil
}

func (s *Server) handlePost(c *routing.Context) error {
	id := c.Param("auth")
	if !s.auth.isAuthenticated(id) {
		return s.redirect(c)
	}

	p, err := promo.FromForm(id, c)
	if err != nil {
		return err
	}

	s.auth.dropAuthentication(id) // disallow further posts for this session
	s.auth.setRateLimited(id)     // mark user as rate limited (can't post again for 30 mins)

	if result, err := s.submit(p); err != nil {
		return err
	} else {
		c.Redirect(result.URL, redirectCode)
	}

	return nil
}

func (s *Server) submit(pr promo.Promotion) (github.PRResponse, error) {
	var empty github.PRResponse

	branch, err0 := s.repo.CreateBranch(pr.ID)
	if err0 != nil {
		return empty, err0
	}

	buf := bytes.Buffer{}
	if err1 := utils.EncodeJson(pr, &buf); err1 != nil {
		return empty, err1
	}

	filename := fmt.Sprintf("%s-%s.json", pr.Type, pr.ID)
	content := buf.Bytes()

	if err3 := branch.CreateFile(filename, content); err3 != nil {
		return empty, err3
	}

	title := fmt.Sprint("Promo for ", pr.Name)
	body := "This PR has been created by a pr-bot!"

	return branch.CreatePR(title, body)
}