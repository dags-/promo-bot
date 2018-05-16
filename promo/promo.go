package promo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dags-/ping/status"
	"github.com/dags-/promo-bot/util"
	"github.com/qiangxue/fasthttp-routing"
)

var (
	discordMatcher = regexp.MustCompile(`(https?://)?(www\.)?(discord\.(gg|io|me|li)|discordapp\.com/invite)/.+[a-z]`)
	websiteMatcher = regexp.MustCompile(`^((http://)|(https://)).*?`)
	ipMatcher      = regexp.MustCompile(`^[a-zA-Z0-9:\\.]+$`)
	urlHint        = "%s url must begin with 'http://' or 'https://'"
	ipHint         = "%s contains invalid characters"
	pingUrl        = "https://ping.dags.me/%s/%s"
)

type Promotion struct {
	ID            string  `json:"id"`
	Type          string  `json:"type"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Icon          string  `json:"icon"`
	Image         string  `json:"image"`
	Website       string  `json:"website"`
	Discord       string  `json:"discord"`
	Tags          string  `json:"tags"`
	IP            *string `json:"ip,omitempty"`
	Whitelist     *bool   `json:"whitelist,omitempty"`
	ServerVersion string  `json:"server_version"`
	ModVersion    string  `json:"mod_version"`
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
	pr.ServerVersion = "unknown"
	pr.ModVersion = "unknown"
	return pr, Validate(&pr)
}

func Validate(pr *Promotion) (error) {
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

	if pr.Type == "server" {
		if pr.IP == nil {
			return errors.New("ip address is required")
		}

		if e := checkValid("ip", *pr.IP, ipHint, ipMatcher, 120); e != nil {
			return e
		}

		if s, m, e := getVersionInfo(*pr.IP); e == nil {
			pr.ServerVersion = s
			pr.ModVersion = m
		} else {
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

func getVersionInfo(address string) (string, string, error) {
	ip, port := splitAddress(address)
	url := fmt.Sprintf(pingUrl, ip, port)

	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	resp, err := client.Get(url)
	if err != nil {
		return "unknown", "unknown", errors.New("unable to connect to server")
	}
	defer resp.Body.Close()

	var s status.Status
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil || s.Type == "error" {
		return "unknown", "unknown", errors.New("unable to read server status")
	}

	if s.Data.ModInfo == nil {
		return s.Data.Version.Name, "unknown", nil
	}

	for _, m := range s.Data.ModInfo.ModList {
		if m.ModID == "conquest" {
			return s.Data.Version.Name, m.Version, nil
		}
	}

	return s.Data.Version.Name, "unknown", nil
}

func splitAddress(address string) (string, string) {
	if strings.ContainsRune(address, ':') {
		split := strings.Split(address, ":")
		if len(split) == 2 {
			return split[0], split[1]
		}
		return split[0], "25565"
	}
	return address, "25565"
}
