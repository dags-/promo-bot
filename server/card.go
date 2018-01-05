package server

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/qiangxue/fasthttp-routing"
	"html/template"
)

var (
	server  = template.Must(template.ParseFiles("_template/card/server.html"))
	twitch  = template.Must(template.ParseFiles("_template/card/twitch.html"))
	youtube = template.Must(template.ParseFiles("_template/card/youtube.html"))
)

func (s *Server) handleSV(c *routing.Context) error {
	return handleCard("servers", s, c)
}

func (s *Server) handleTW(c *routing.Context) error {
	return handleCard("twitchers", s, c)
}

func (s *Server) handleYT(c *routing.Context) error {
	return handleCard("youtubers", s, c)
}

func handleCard(group string, s *Server, c *routing.Context) error {
	promoId := c.Param("id")
	promo, err := s.api.GetPromo(group, promoId)
	if err != nil {
		return err
	}

	temp, err := getTemplate(group)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	err = temp.Execute(&buf, promo)
	if err != nil {
		return err
	}

	c.Response.Header.Set("Content-Type", "text/html")
	c.Response.SetBody(buf.Bytes())
	return nil
}

func getTemplate(group string) (*template.Template, error) {
	switch group {
	case "servers":
		return server, nil
	case "twitchers":
		return twitch, nil
	case "youtubers":
		return youtube, nil
	default:
		return nil, errors.New("unknown group: " + group)
	}
}
