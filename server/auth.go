package server

import (
	"github.com/qiangxue/fasthttp-routing"
	"net/http"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"errors"
)

const authUrl = "https://discordapp.com/oauth2/authorize"
const userUrl = "https://discordapp.com/api/users/@me"
const tokenUrl = "https://discordapp.com/api/oauth2/token"

func (s *Server) redirect(c *routing.Context) error {
	form := url.Values{}
	form.Set("scope", "identify")
	form.Set("client_id", s.clientId)
	form.Set("response_type", "code")
	form.Set("redirect_uri", s.redirectUri)
	u := fmt.Sprint(authUrl, "?", form.Encode())
	c.Redirect(u, 302)
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
	dec := json.NewDecoder(resp.Body)
	err2 := dec.Decode(&tokenAuth)
	if err2 != nil {
		return err2
	}

	token := tokenAuth["access_token"]
	req, _ = http.NewRequest("GET", userUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", "handleLogin (n/a, 1.0)")
	resp, err3 := client.Do(req)
	if err3 != nil {
		return err3
	}

	var user map[string]interface{}
	dec = json.NewDecoder(resp.Body)
	err4 := dec.Decode(&user)
	if err4 != nil {
		return err4
	}

	raw, ok := user["id"];
	if !ok {
		return errors.New("Could not authenticate")
	}

	auth := fmt.Sprint(raw)
	redirect := fmt.Sprint(promotionRoute, "/", auth)
	s.auth.setAuthenticated(auth)
	c.Redirect(redirect, 302)
	return nil
}