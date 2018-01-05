package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dags-/promo-bot/promo"
	"time"
)

func (b *Bot) StartLoop() {
	go func() {
		pause := time.Duration(15) * time.Second
		for {
			promos := b.api.GetPromoQueue()

			for _, pr := range promos {
				uid := fmt.Sprint(pr.GetMeta().Type, "-", pr.GetMeta().ID)

				if uid != b.lastPromo {
					b.postPromotion(pr)
					b.lastPromo = uid
					time.Sleep(b.config.getInterval())
				}
			}

			time.Sleep(pause)
		}
	}()
}

func (b *Bot) postPromotion(pr promo.Promo) {
	meta := pr.GetMeta()

	embed := &discordgo.MessageEmbed{
		Title:       meta.Name,
		URL:         meta.Website,
		Description: meta.Description,
		Author: &discordgo.MessageEmbedAuthor{
			URL: meta.Discord,
			IconURL: meta.Icon,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: meta.Icon,
		},
	}

	switch meta.Type {
	case "server":
		s := pr.(*promo.Server)
		embed.Color = 0x00d56a
		embed.Author.Name = "#Server"
		embed.Fields = []*discordgo.MessageEmbedField{
			{Name: "IP", Value: s.IP, Inline: true},
			{Name: "Whitelist", Value: fmt.Sprint(s.Whitelist), Inline: true},
		}
		break
	case "twitcher":
		embed.Color = 0x0080ff
		embed.Author.Name = "#Twitcher"
		break
	case "youtuber":
		embed.Color = 0xff8080
		embed.Author.Name = "#Youtuber"
		break
	}

	// video embed not working :/
	if meta.Media.Type == "image" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: meta.Media.URL,
		}
	}

	message := &discordgo.MessageSend{
		Content: meta.Discord,
		Embed:   embed,
	}

	channels := b.config.getChannels()
	for _, ch := range channels {
		if _, err := b.sess.Channel(ch); err != nil {
			fmt.Println("Channel doesn't exist: ", ch)
			b.config.removeChannel(ch)
			continue
		}

		_, err := b.sess.ChannelMessageSendComplex(ch, message)
		if err != nil {
			fmt.Println("Post err: ", err)
		}
	}
}
