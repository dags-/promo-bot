package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dags-/promo-bot/promo"
	"os"
)

var (
	owner = "99824915045191680"
)

type Bot struct {
	config    *Config
	lastPromo string
	sess      *discordgo.Session
	api       *promo.Api
}

func NewBot(a *promo.Api) *Bot {
	return &Bot{
		api:    a,
		config: getOrCreate(),
	}
}

func (b *Bot) Start(token string) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
		return
	}

	b.sess = s
	b.sess.AddHandler(b.command)
	b.sess.AddHandler(b.ready)
	b.sess.AddHandler(b.join)

	err = b.sess.Open()
	defer b.sess.Close()
	if err != nil {
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	<-sc
}