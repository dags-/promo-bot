package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dags-/promo-bot/promo"
	"github.com/dags-/promo-bot/server"
	"os"
	"time"
)

var (
	owner = "99824915045191680"
)

type Bot struct {
	config    Config
	lastPromo string
	sess      *discordgo.Session
	api       *server.Api
}

func NewBot(a *server.Api) *Bot {
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

func (b *Bot) StartLoop() {
	go func() {
		for {
			promos := b.api.GetPromoQueue()

			for _, pr := range promos {
				uid := fmt.Sprint(pr.GetMeta().Type, "-", pr.GetMeta().ID)

				if uid == b.lastPromo {
					continue
				}

				b.post(pr)
				b.lastPromo = uid
				time.Sleep(b.config.getInterval())
			}
		}
	}()
}

func (b *Bot) post(pr promo.Promo) {
	meta := pr.GetMeta()
	message := &discordgo.MessageEmbed{
		Title:       meta.Name,
		URL:         meta.Link,
		Description: meta.Description,
		Author: &discordgo.MessageEmbedAuthor{
			URL:     meta.Link,
			IconURL: meta.Icon,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: meta.Icon,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: meta.Media,
		},
	}

	switch meta.Type {
	case "server":
		s := pr.(*promo.Server)
		message.Color = 0x00d56a
		message.Author.Name = "#Server"
		message.Fields = []*discordgo.MessageEmbedField{
			{Name: "IP", Value: s.IP, Inline: true},
			{Name: "Whitelist", Value: fmt.Sprint(s.Whitelist), Inline: true},
		}
		break
	case "twitcher":
		message.Color = 0x0080ff
		message.Author.Name = "#Twitcher"
		break
	case "youtuber":
		message.Color = 0xff8080
		message.Author.Name = "#Youtuber"
		break
	}

	channels := b.config.getChannels()
	for _, ch := range channels {
		b.sess.ChannelMessageSendEmbed(ch, message)
	}
}
