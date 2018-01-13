package promo

import (
	"errors"
	"fmt"
	"github.com/dags-/promo-bot/util"
	"github.com/qiangxue/fasthttp-routing"
	"regexp"
)

var (
	discordMatcher = regexp.MustCompile(`(https?://)?(www\.)?(discord\.(gg|io|me|li)|discordapp\.com/invite)/.+[a-z]`)
	websiteMatcher = regexp.MustCompile(`^((http://)|(https://)).*?`)
	ipMatcher      = regexp.MustCompile(`^[a-zA-Z0-9:\\.]+$`)
	urlHint        = "%s url must begin with 'http://' or 'https://'"
	ipHint         = "%s contains invalid characters"
)

type Promotion struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	Image       string  `json:"image"`
	Website     string  `json:"website"`
	Discord     string  `json:"discord"`
	Tags        string  `json:"tags"`
	IP          *string `json:"ip,omitempty"`
	Whitelist   *bool   `json:"whitelist,omitempty"`
}

func FromForm(id string, c *routing.Context) (Promotion, error) {
	var pr Promotion
	pr.ID = id
	pr.Type = utils.String(c, "type")
	pr.Name = utils.String(c, "name")
	pr.Description = utils.String(c, "description")
	pr.Icon = utils.String(c, "icon")
	pr.Image = utils.String(c, "image")
	pr.Website = utils.String(c, "website")
	pr.Discord = utils.String(c, "discord")
	pr.Tags = utils.String(c, "tags")
	pr.IP = utils.StringOp(c, "ip")
	pr.Whitelist = utils.BoolOp(c, "whitelist")
	return pr, Validate(pr)
}

func Validate(pr Promotion) (error) {
	if pr.ID == "" {
		return errors.New("id is missing")
	}

	if pr.Name == "" {
		return errors.New("name is required")
	}

	if len(pr.Name) > 120 {
		return errors.New("name is too long")
	}

	if len(pr.Description) > 480 {
		return errors.New("description is too long")
	}

	if e := checkValid("icon", pr.Icon, urlHint, websiteMatcher, 240); e != nil {
		return e
	}

	if e := checkValid("image", pr.Image, urlHint, websiteMatcher, 240); e != nil {
		return e
	}

	if e := checkValid("website", pr.Website, urlHint, websiteMatcher, 240); e != nil {
		return e
	}

	if e := checkValid("discord", pr.Discord, urlHint, discordMatcher, 120); e != nil {
		return e
	}

	if pr.IP != nil {
		if e := checkValid("ip", *pr.IP, ipHint, ipMatcher, 120); e != nil {
			return e
		}
	}

	return nil
}

func checkValid(name, url, hint string, match *regexp.Regexp, max int) (error) {
	if url == "" {
		return nil
	}

	if len(url) > max {
		return fmt.Errorf("%s is too long", name)
	}

	if !match.MatchString(url) {
		return fmt.Errorf(hint, name)
	}

	return nil
}
