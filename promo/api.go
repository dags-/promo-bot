package promo

import (
	"fmt"
	"github.com/dags-/promo-bot/github"
	"github.com/dags-/promo-bot/util"
	"github.com/pkg/errors"
	"github.com/qiangxue/fasthttp-routing"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

type Api struct {
	lock sync.RWMutex
	repo github.Repo
	data map[string]map[string]Promotion
}

func NewApi(repo github.Repo) *Api {
	return &Api{
		repo: repo,
	}
}

func (api *Api) HandleGet(c *routing.Context) error {
	api.lock.RLock()
	defer api.lock.RUnlock()

	var response interface{}

	promoType := c.Param("type")
	promoId := c.Param("id")
	if promoType == "all" {
		response = api.data
	} else {
		data, ok := api.data[promoType]
		if ok {
			response = data
			if promoId != "" {
				val, ok := data[promoId]
				if ok {
					response = val
				}
			}
		}
	}

	if response == nil {
		return errors.New(fmt.Sprint("unknown route, type: ", promoType, ", id: ", promoId))
	}


	c.Response.Header.SetStatusCode(http.StatusOK)
	c.Response.Header.Set("Cache-Control", "max-age=300")
	c.Response.Header.SetContentType("application/json; charset=UTF-8")

	return utils.EncodeJson(response, c.Response.BodyWriter())
}

func (api *Api) GetPromoQueue() ([]Promotion) {
	api.lock.RLock()
	defer api.lock.RUnlock()

	var promos []Promotion
	for _, ps := range api.data {
		for _, p := range ps {
			promos = append(promos, p)
		}
	}

	queue := make([]Promotion, len(promos))
	perm := rand.Perm(len(promos))
	for i, v := range perm {
		queue[v] = promos[i]
	}

	return queue
}

func (api *Api) Tick() {
	fmt.Println("Updating api...")
	api.lock.Lock()
	defer api.lock.Unlock()

	contents, err := api.repo.GetContents()
	if err != nil {
		fmt.Println(err)
		return
	}

	data := map[string]map[string]Promotion {
		"server": make(map[string]Promotion),
		"twitch": make(map[string]Promotion),
		"youtube": make(map[string]Promotion),
	}

	for _, c := range contents {
		if !strings.HasSuffix(c.Name, ".json") {
			continue
		}

		resp, err := http.Get(c.URL)
		if err != nil {
			continue
		}

		var pr Promotion
		if err := utils.DecodeJson(&pr, resp.Body); err != nil {
			fmt.Println("Err api.tick.decode: ", err)
			continue
		}

		if promos, ok := data[pr.Type]; ok {
			promos[pr.ID] = pr
		} else {
			fmt.Println("Err api.tick.data: invalid promo type: ", pr.Type)
		}
	}

	api.data = data
	fmt.Println("Remaining requests: ", )
}
