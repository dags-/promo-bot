package server

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/dags-/promo-bot/github"
	"github.com/valyala/fasthttp"
	"time"
	"fmt"
	"io"
	"encoding/json"
)

const apiRoute = "/api/<type>"
const apiIdRoute = "/api/<type>/<id>"
const previewRoute = "/preview/<type>/<id>"
const authRoute = "/auth"
const promotionRoute = "/promotion"
const promotionAuthRoute = "/promotion/<auth>"

type Server struct {
	session      github.Session
	repo         github.Repo
	auth         AuthSessions
	api          Api
	clientId     string
	clientSecret string
	redirectUri  string
}

func NewServer(s github.Session, r github.Repo, clientId, clientSecret, redirectUri string) Server {
	return Server{
		session: s,
		repo: r,
		clientId: clientId,
		clientSecret: clientSecret,
		redirectUri: redirectUri,
		api: newApi(r),
		auth: newAuthSessions(),
	}
}

func (s *Server) Start(port int) {
	router := routing.New()
	router.Get(apiRoute, s.handleApi)
	router.Get(apiIdRoute, s.handleApi)
	router.Get(previewRoute, s.handlePreview)
	router.Get(authRoute, s.handleAuth)
	router.Get(promotionRoute, s.redirect)
	router.Get(promotionAuthRoute, s.handleAppGet)
	router.Post(promotionAuthRoute, s.handleAppPost)

	server := fasthttp.Server{
		Handler: router.HandleRequest,
		GetOnly: false,
		DisableKeepalive: true,
		ReadBufferSize: 10240,
		WriteBufferSize: 25600,
		ReadTimeout: time.Duration(time.Second * 2),
		WriteTimeout: time.Duration(time.Second * 2),
		MaxConnsPerIP: 3,
		MaxRequestsPerConn: 1,
		MaxRequestBodySize: 0,
	}

	go startServerLoop(s)

	panic(server.ListenAndServe(fmt.Sprintf(":%v", port)))
}

func startServerLoop(s *Server) {
	sleep := time.Duration(time.Minute * 15)
	for {
		s.api.tick()
		s.auth.tick()
		time.Sleep(sleep)
	}
}

func getString(c *routing.Context, key string) (string) {
	d := c.FormValue(key)
	if d != nil {
		return string(d)
	}
	return ""
}

func toJson(i interface{}, wr io.Writer) (error) {
	enc := json.NewEncoder(wr)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(i)
}