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
	Link        string `json:"link"`
	Media       string `json:"media"`
}

type Server struct {
	Meta
	IP        string `json:"ip"`
	Whitelist bool   `json:"whitelist"`
}

type Youtuber struct {
	Meta
}

type Twitcher struct {
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
	case "youtuber":
		var y Youtuber
		return &y, json.Unmarshal(data, &y)
	case "twitcher":
		var t Twitcher
		return &t, json.Unmarshal(data, &t)
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

func (t Twitcher) GetMeta() (*Meta) {
	return &t.Meta
}

func (y Youtuber) GetMeta() (*Meta) {
	return &y.Meta
}
