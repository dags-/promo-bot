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

var random = rand.New(rand.NewSource(8008))

type Api struct {
	lock    sync.RWMutex
	repo    github.Repo
	Server  map[string]promo.Promo `json:"server"`
	Twitch  map[string]promo.Promo `json:"twitch"`
	Youtube map[string]promo.Promo `json:"youtube"`
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

	var promos []promo.Promo
	for _, p := range api.Server {
		promos = append(promos, p)
	}
	for _, p := range api.Twitch {
		promos = append(promos, p)
	}
	for _, p := range api.Youtube {
		promos = append(promos, p)
	}

	if len(promos) > 1 {
		return promos
	}

	queue := make([]promo.Promo, len(promos))
	perm := random.Perm(len(promos))
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
	case "server":
		return api.Server, nil
	case "twitch":
		return api.Twitch, nil
	case "youtube":
		return api.Youtube, nil
	default:
		return nil, errors.New("invalid promo type")
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

	server := make(map[string]promo.Promo)
	twitch := make(map[string]promo.Promo)
	youtube := make(map[string]promo.Promo)

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
			server[pr.GetMeta().ID] = pr
			break
		case "youtube":
			youtube[pr.GetMeta().ID] = pr
			break
		case "twitch":
			twitch[pr.GetMeta().ID] = pr
			break
		}
	}

	api.Server = server
	api.Youtube = youtube
	api.Twitch = twitch
}
