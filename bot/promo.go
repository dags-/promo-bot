package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dags-/promo-bot/promo"
	"time"
)

func (b *Bot) StartLoop() {
	go func() {
		h := NewHistory(5)

		for {
			promos := b.api.GetPromoQueue()
			start := time.Now()

			fmt.Printf("Starting new promotions run: %v\n", len(promos))
			for _, pr := range promos {
				if h.Contains(pr.GetMeta().ID) {
					continue
				}

				b.postPromotion(pr)
				h.Add(pr.GetMeta().ID)
				time.Sleep(b.config.getInterval())
			}

			remaining := b.config.getInterval() - time.Since(start)
			if remaining > 0 {
				time.Sleep(remaining)
			}
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
			URL: promo.Or(meta.Discord == "", meta.Website, meta.Discord),
			IconURL: meta.Icon,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: meta.Icon,
		},
		Provider: &discordgo.MessageEmbedProvider{
			Name: "promo-bot",
		},
	}

	switch meta.Type {
	case "server":
		s := pr.(*promo.Server)
		embed.Color = 0x00d56a
		embed.Author.Name = "#Server"
		embed.Fields = []*discordgo.MessageEmbedField{
			{Name: "IP", Value: s.IP, Inline: true},
			{Name: "Whitelist", Value: promo.Or(s.Whitelist, "Yes", "No"), Inline: true},
		}
		break
	case "twitch":
		embed.Color = 0x0080ff
		embed.Author.Name = "#Twitch"
		break
	case "youtube":
		embed.Color = 0xff8080
		embed.Author.Name = "#Youtube"
		break
	}

	message := &discordgo.MessageSend{
		Embed: embed,
	}

	if meta.Media.Type == "image" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: meta.Media.URL,
		}
	} else if meta.Media.Type == "video" { // video embed not working :/
		embed.Video = &discordgo.MessageEmbedVideo{
			URL: meta.Media.URL,
		}
	}

	channels := b.config.getChannels()
	for _, ch := range channels {
		if _, err := b.sess.Channel(ch); err != nil {
			fmt.Println("Channel doesn't exist: ", ch)
			b.config.removeChannel(ch)
			continue
		}

		_, err := b.sess.ChannelMessageSendComplex(ch, message)
		if err == nil && meta.Media.Type == "video" {
			b.sess.ChannelMessageSend(ch, meta.Media.URL)
		}

		if err != nil {
			fmt.Println("Post err: ", err)
		}
	}
}