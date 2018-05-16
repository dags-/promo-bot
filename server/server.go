package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/dags-/promo-bot/github"
	"github.com/dags-/promo-bot/promo"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

const redirectCode = 301
const filesRoute = "/files/*"
const addBotRoute = "/bot"
const apiRoute = "/api/<type>"
const apiIdRoute = "/api/<type>/<id>"
const authRoute = "/auth"
const promotionRoute = "/apply"
const promotionAuthRoute = "/apply/<auth>"

type Server struct {
	session      github.Session
	repo         github.Repo
	auth         AuthSessions
	Api          *promo.Api
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
		Api:          promo.NewApi(r),
		auth:         newAuthSessions(),
	}
}

func (s *Server) Start(port int) {
	router := routing.New()
	router.Get(apiRoute, s.Api.HandleGet)
	router.Get(apiIdRoute, s.Api.HandleGet)
	router.Get(authRoute, s.handleAuth)
	router.Get(promotionRoute, s.redirect)
	router.Get(promotionAuthRoute, s.handleGet)
	router.Post(promotionAuthRoute, s.handlePost)
	router.Get(addBotRoute, s.handleAdd)
	router.Get(filesRoute, newFileHandler())

	router.Use()

	server := fasthttp.Server{
		Handler:            router.HandleRequest,
		GetOnly:            false,
		DisableKeepalive:   true,
		ReadBufferSize:     10240,
		WriteBufferSize:    25600,
		ReadTimeout:        time.Duration(time.Second * 2),
		WriteTimeout:       time.Duration(time.Second * 2),
		MaxConnsPerIP:      3,
		MaxRequestsPerConn: 1,
		MaxRequestBodySize: 0,
	}

	go startServerLoop(s)

	panic(server.ListenAndServe(fmt.Sprintf(":%v", port)))
}

func newFileHandler() (func(context *routing.Context) error) {
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
	sleep := time.Duration(time.Minute * 10)
	for {
		s.Api.Tick()
		s.auth.tick()
		go func() {
			r, err := s.session.RemainingRate()
			if err == nil {
				fmt.Println("Remaining github calls: ", r)
			}
		}()
		time.Sleep(sleep)
	}
}
