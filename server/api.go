package server

import (
	"errors"
	"fmt"
	"github.com/dags-/promo-bot/github"
	"github.com/dags-/promo-bot/promo"
	"github.com/qiangxue/fasthttp-routing"
	"math/rand"
	"net/http"
	"sync"
)

type Api struct {
	lock      sync.RWMutex
	repo      github.Repo
	Servers   map[string]promo.Promo `json:"servers"`
	Youtubers map[string]promo.Promo `json:"youtubers"`
	Twitchers map[string]promo.Promo `json:"twitchers"`
}

func newApi(repo github.Repo) *Api {
	return &Api{
		repo: repo,
	}
}

func (s *Server) handleApi(c *routing.Context) error {
	api := s.Api
	promoType := c.Param("type")
	promoId := c.Param("id")

	if promoType == "all" {
		api.lock.Lock()
		defer api.lock.Unlock()
		return toJson(api, c.Response.BodyWriter())
	}

	promos, err := api.GetType(promoType)
	if err != nil {
		return err
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

func (api *Api) GetPromoQueue() ([]promo.Promo) {
	api.lock.Lock()
	defer api.lock.Unlock()

	promos := []promo.Promo{}
	for _, p := range api.Servers {
		promos = append(promos, p)
	}
	for _, p := range api.Twitchers {
		promos = append(promos, p)
	}
	for _, p := range api.Youtubers {
		promos = append(promos, p)
	}

	if len(promos) > 1 {
		return promos
	}

	queue := make([]promo.Promo, len(promos))
	perm := rand.Perm(len(promos))
	for i, v := range perm {
		queue[v] = promos[i]
	}

	return queue
}

func (api *Api) GetPromo(promoType, promoId string) (promo.Promo, error) {
	promos, err := api.GetType(promoType)
	if err != nil {
		return nil, err
	}

	if pr, ok := promos[promoId]; ok {
		return pr, nil
	}

	return nil, errors.New("promotion for id not found")
}

func (api *Api) GetType(promoType string) (map[string]promo.Promo, error) {
	api.lock.Lock()
	defer api.lock.Unlock()

	switch promoType {
	case "servers":
		return api.Servers, nil
	case "youtubers":
		return api.Youtubers, nil
	case "twitchers":
		return api.Twitchers, nil
	default:
		return nil, errors.New("Invalid promo type")
	}
}

func (api *Api) tick() {
	api.lock.Lock()
	defer api.lock.Unlock()

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
		case "youtube":
			youtubers[pr.GetMeta().ID] = pr
			break
		case "twitch":
			twitchers[pr.GetMeta().ID] = pr
			break
		}
	}

	api.Servers = servers
	api.Youtubers = youtubers
	api.Twitchers = twitchers
}
