package promo

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"errors"
)

type Promo interface {
	GetMeta() *Meta
}

type Meta struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Media       string `json:"media"`
	Discord     string `json:"discord"`
}

type Server struct {
	Meta
	IP        string `json:"ip"`
	Whitelist bool `json:"whitelist"`
	Website   string `json:"website"`
}

type Youtuber struct {
	Meta
	ChannelName string `json:"title"`
	URL         string `json:"url"`
}

type Twitcher struct {
	Meta
	UserName string `json:"username"`
	URL      string `json:"url"`
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