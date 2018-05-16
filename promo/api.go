package promo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"

	"github.com/dags-/promo-bot/github"
	"github.com/dags-/promo-bot/util"
	"github.com/pkg/errors"
	"github.com/qiangxue/fasthttp-routing"
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
	} else if data, ok := api.data[promoType]; ok {
		if promoId == "" {
			response = data
		} else if val, ok := data[promoId]; ok {
			response = val
		} else {
			return errors.Errorf("no %s promo for id %s", promoType, promoId)
		}
	} else {
		return errors.Errorf("unknown promo type %s", promoType)
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
	promos = append(promos, loadSelf())

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
	promos, err := GetPromotions(api.repo.Owner, api.repo.Name)

	if err != nil {
		fmt.Println("Error fetching promotions: ", err)
		return
	}

	api.lock.Lock()
	defer api.lock.Unlock()
	api.data = promos
}

func loadSelf() (Promotion) {
	var self Promotion
	data, err := ioutil.ReadFile("self.json")
	if err != nil {
		return self
	}
	json.Unmarshal(data, &self)
	return self
}
