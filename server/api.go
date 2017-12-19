package server

import (
	"github.com/dags-/promo-bot/promo"
	"sync"
	"fmt"
	"net/http"
	"github.com/qiangxue/fasthttp-routing"
	"errors"
	"github.com/dags-/promo-bot/github"
)

type Api struct {
	lock      sync.RWMutex
	repo      github.Repo
	Servers   map[string]promo.Promo `json:"servers"`
	Youtubers map[string]promo.Promo `json:"youtubers"`
	Twitchers map[string]promo.Promo `json:"twitchers"`
}

func newApi(repo github.Repo) Api {
	return Api{
		repo: repo,
	}
}

func (s *Server) handleApi(c *routing.Context) error {
	s.api.lock.Lock()
	defer s.api.lock.Unlock()

	promoType := c.Param("type")
	promoId := c.Param("id")

	var promos map[string]promo.Promo

	switch promoType {
	case "servers":
		promos = s.api.Servers
		break
	case "youtubers":
		promos = s.api.Youtubers
		break
	case "twitchers":
		promos = s.api.Twitchers
		break
	case "all":
		return toJson(s.api, c.Response.BodyWriter())
	default:
		return errors.New("Invalid roote, try `/servers` `/youtubers` `/twitchers`")
	}

	if promoId != "" {
		if p, ok := promos[promoId]; ok {
			return toJson(p, c.Response.BodyWriter())
		} else {
			err := fmt.Sprintf("No <%s> promotion for id id <%s>", promoType, promoId)
			return errors.New(err)
		}
	}

	return toJson(promos, c.Response.BodyWriter())
}

func (api *Api) tick() {
	contents, err := api.repo.GetContents()
	if err != nil {
		fmt.Println(err)
		return
	}

	servers := make(map[string]promo.Promo)
	youtubers := make(map[string]promo.Promo)
	twitchers := make(map[string]promo.Promo)

	for _, c := range contents {
		resp, err := http.Get(c.URL)
		if err != nil {
			continue
		}

		pr, err := promo.Read(resp.Body)
		if err != nil {
			continue
		}

		switch pr.GetMeta().Type {
		case "server":
			servers[pr.GetMeta().ID] = pr
			break
		case "youtuber":
			youtubers[pr.GetMeta().ID] = pr
			break
		case "twitcher":
			twitchers[pr.GetMeta().ID] = pr
			break
		}
	}

	api.lock.Lock()
	defer api.lock.Unlock()
	api.Servers = servers
	api.Youtubers = youtubers
	api.Twitchers = twitchers
}