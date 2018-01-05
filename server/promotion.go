package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dags-/promo-bot/github"
	"github.com/dags-/promo-bot/promo"
	"github.com/qiangxue/fasthttp-routing"
	"regexp"
)

var (
	discordmatcher = regexp.MustCompile(`(https?://)?(www\.)?(discord\.(gg|io|me|li)|discordapp\.com/invite)/.+[a-z]`)
)

func (s *Server) handleAppGet(c *routing.Context) error {
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

func (s *Server) handleAppPost(c *routing.Context) error {
	id := c.Param("auth")
	if !s.auth.isAuthenticated(id) {
		return s.redirect(c)
	}

	var p promo.Promo

	meta := promo.Meta{
		ID:          id,
		Type:        getString(c, "type"),
		Name:        getString(c, "name"),
		Website:     getString(c, "website"),
		Discord:     getString(c, "discord"),
		Icon:        getString(c, "icon"),
		Description: getString(c, "description"),
		Media: promo.Media{
			Type: getString(c, "media-type"),
			URL:  getString(c, "media-url"),
		},
	}

	err := validate(meta)
	if err != nil {
		return err
	}

	s.auth.dropAuthentication(id) // disallow further posts for this session
	s.auth.setRateLimited(id)     // mark user as rate limited (can't post again for 30 mins)

	switch meta.Type {
	case "server":
		wl := getString(c, "whitelist") != ""
		var server promo.Server
		server.Meta = meta
		server.IP = getString(c, "ip")
		server.Website = getString(c, "link")
		server.Whitelist = wl
		p = &server
		break
	case "youtuber":
		var youtuber promo.Youtuber
		youtuber.Meta = meta
		p = &youtuber
		break
	case "twitcher":
		var twitcher promo.Twitcher
		twitcher.Meta = meta
		p = &twitcher
		break
	}

	if p.GetMeta().ID == "" {
		return errors.New("no id associated with the promotion")
	}

	result, err := s.submit(p)
	if err != nil {
		return err
	}

	c.Redirect(result.URL, 302)

	return nil
}

func validate(meta promo.Meta) (error) {
	if meta.Name == "" {
		return errors.New("name is required")
	}

	if len(meta.Name) > 120 {
		return errors.New("name is too long")
	}

	if len(meta.Description) > 240 {
		return errors.New("description is too long")
	}

	if meta.Discord != "" && !discordmatcher.MatchString(meta.Discord) {
		return errors.New("invalid discord link")
	}

	return nil
}

func (s *Server) submit(promo promo.Promo) (github.PRResponse, error) {
	var empty github.PRResponse

	branch, err0 := s.repo.CreateBranch(promo.GetMeta().ID)
	if err0 != nil {
		return empty, err0
	}

	buf := bytes.Buffer{}
	err1 := toJson(promo, &buf)
	if err1 != nil {
		return empty, err1
	}

	filename := fmt.Sprintf("%s-%s.json", promo.GetMeta().Type, promo.GetMeta().ID)
	content := buf.Bytes()

	err3 := branch.CreateFile(filename, content)
	if err3 != nil {
		return empty, err3
	}

	title := fmt.Sprint("Promo for ", promo.GetMeta().Name)
	body := "This PR has been created by a promo-bot!"

	return branch.CreatePR(title, body)
}
