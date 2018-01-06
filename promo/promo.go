package promo

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

type Promo interface {
	GetMeta() *Meta
}

type Meta struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Website     string `json:"website"`
	Discord     string `json:"discord"`
	Media       Media  `json:"media"`
}

type Media struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Server struct {
	Meta
	IP        string `json:"ip"`
	Whitelist bool   `json:"whitelist"`
}

type Youtube struct {
	Meta
}

type Twitch struct {
	Meta
}

func Read(r io.Reader) (Promo, error) {
	var meta Meta
	data, err := ioutil.ReadAll(r)

	if err != nil {
		return &meta, err
	}

	err = json.Unmarshal(data, &meta)
	if err != nil {
		return &meta, err
	}

	switch meta.Type {
	case "server":
		var s Server
		return &s, json.Unmarshal(data, &s)
	case "twitch":
		var t Twitch
		return &t, json.Unmarshal(data, &t)
	case "youtube":
		var y Youtube
		return &y, json.Unmarshal(data, &y)
	default:
		return &meta, errors.New("Invalid promo type: " + meta.Type)
	}
}

func (m *Meta) GetMeta() (*Meta) {
	return m
}

func (s Server) GetMeta() (*Meta) {
	return &s.Meta
}

func (t Twitch) GetMeta() (*Meta) {
	return &t.Meta
}

func (y Youtube) GetMeta() (*Meta) {
	return &y.Meta
}

func Or(exp bool, a, b string) string {
	if exp {
		return a
	}
	return b
}