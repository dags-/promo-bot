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
			interval := b.config.getInterval()

			fmt.Printf("Starting new promotions run: %v\n", len(promos))
			for _, pr := range promos {
				uid := fmt.Sprint(pr.GetMeta().Type, "-", pr.GetMeta().ID)

				if h.Contains(uid) {
					continue
				}

				b.postPromotion(pr)
				h.Add(uid)
				time.Sleep(interval)
			}

			remaining := interval - time.Since(start)
			if remaining > 0 {
				fmt.Println("Sleeping remaining time: ", remaining)
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
			{Name: "Whitelist", Value: promo.Or(s.Whitelist, "Yes", "No"), Inline: true},
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