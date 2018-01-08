package promo

import (
	"errors"
	"github.com/dags-/promo-bot/util"
	"github.com/qiangxue/fasthttp-routing"
	"regexp"
)

var (
	discordMatcher = regexp.MustCompile(`(https?://)?(www\.)?(discord\.(gg|io|me|li)|discordapp\.com/invite)/.+[a-z]`)
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

	if len(pr.Icon) > 240 {
		return errors.New("icon url too long")
	}

	if len(pr.Image) > 240 {
		return errors.New("image url too long")
	}

	if len(pr.Website) > 240 {
		return errors.New("website url too long")
	}

	if pr.Discord != "" && !discordMatcher.MatchString(pr.Discord) {
		return errors.New("invalid discord link")
	}

	if pr.IP != nil && len(*pr.IP) > 120 {
		return errors.New("ip address too long")
	}

	return nil
}
