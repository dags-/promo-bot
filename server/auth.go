package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dags-/promo-bot/util"
	"github.com/qiangxue/fasthttp-routing"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const authUrl = "https://discordapp.com/oauth2/authorize"
const userUrl = "https://discordapp.com/api/users/@me"
const tokenUrl = "https://discordapp.com/api/oauth2/token"

var (
	home = template.Must(template.ParseFiles("_template/home.html"))
)

func (s *Server) redirect(c *routing.Context) error {
	form := url.Values{}
	form.Set("scope", "identify")
	form.Set("client_id", s.clientId)
	form.Set("response_type", "code")
	form.Set("redirect_uri", s.redirectUri)
	u := fmt.Sprint(authUrl, "?", form.Encode())
	buf := bytes.Buffer{}
	err := home.Execute(&buf, u)
	if err != nil {
		return err
	}
	c.Response.Header.Set("Content-Type", "text/html")
	c.Response.SetBody(buf.Bytes())
	return nil
}

func (s *Server) handleAuth(c *routing.Context) error {
	code := string(c.FormValue("code"))

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_secret", s.clientSecret)
	form.Set("client_id", s.clientId)
	form.Set("redirect_uri", s.redirectUri)
	form.Set("code", code)

	req, err0 := http.NewRequest("POST", tokenUrl, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err0 != nil {
		return err0
	}

	client := http.Client{}
	resp, err1 := client.Do(req)
	if err1 != nil {
		return err1
	}

	var tokenAuth map[string]interface{}
	if err2 := utils.DecodeJson(&tokenAuth, resp.Body); err2 != nil {
		return err2
	}

	token := tokenAuth["access_token"]
	req, _ = http.NewRequest("GET", userUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", "PromoBot (n/a, 1.0)")
	resp, err3 := client.Do(req)
	if err3 != nil {
		return err3
	}

	var user map[string]interface{}
	if err4 := utils.DecodeJson(&user, resp.Body); err4 != nil {
		return err4
	}

	raw, ok := user["id"]
	if !ok {
		return errors.New("could not authenticate")
	}

	id := raw.(string)
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	auth := strconv.FormatInt(i, 36)
	redirect := fmt.Sprint(promotionRoute, "/", auth)
	s.auth.setAuthenticated(auth)
	c.Redirect(redirect, redirectCode)
	return nil
}
