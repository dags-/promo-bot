package server

import (
	"encoding/json"
	"fmt"
	"github.com/dags-/promo-bot/github"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"io"
	"strings"
	"time"
)

const apiRoute = "/Api/<type>"
const apiIdRoute = "/Api/<type>/<id>"
const serverCardRoute = "/sv/<id>"
const twitchCardRoute = "/tw/<id>"
const youtubeCardRoute = "/yt/<id>"
const authRoute = "/auth"
const filesRoute = "/files/*"
const promotionRoute = "/apply"
const promotionAuthRoute = "/apply/<auth>"

type Server struct {
	session      github.Session
	repo         github.Repo
	auth         AuthSessions
	Api          Api
	clientId     string
	clientSecret string
	redirectUri  string
}

type PathMap map[string]string

func NewServer(s github.Session, r github.Repo, clientId, clientSecret, redirectUri string) Server {
	return Server{
		session:      s,
		repo:         r,
		clientId:     clientId,
		clientSecret: clientSecret,
		redirectUri:  redirectUri,
		Api:          newApi(r),
		auth:         newAuthSessions(),
	}
}

func (s *Server) Start(port int) {
	router := routing.New()
	router.Get(apiRoute, s.handleApi)
	router.Get(apiIdRoute, s.handleApi)
	router.Get(serverCardRoute, s.handleSV)
	router.Get(twitchCardRoute, s.handleTW)
	router.Get(youtubeCardRoute, s.handleYT)
	router.Get(authRoute, s.handleAuth)
	router.Get(promotionRoute, s.redirect)
	router.Get(promotionAuthRoute, s.handleAppGet)
	router.Post(promotionAuthRoute, s.handleAppPost)
	router.Get(filesRoute, newFileHandler())

	router.Use()

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

func newFileHandler() (func(context *routing.Context)error) {
	prefix := strings.TrimSuffix(filesRoute, "/*")
	split := len([]byte(prefix))

	fs := fasthttp.FS{
		Root: "_public/",
		PathRewrite: func(ctx *fasthttp.RequestCtx) []byte {
			return ctx.Path()[split:]
		},
	}

	handler := fs.NewRequestHandler()
	return func(c *routing.Context) error {
		handler(c.RequestCtx)
		return nil
	}
}

func startServerLoop(s *Server) {
	sleep := time.Duration(time.Minute * 15)
	for {
		s.Api.tick()
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