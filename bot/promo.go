package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dags-/promo-bot/promo"
	"strings"
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

				m := buildMessage(pr)
				b.sendToAll(m, pr)
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

func (b *Bot) sendToAll(m *discordgo.MessageSend, pr promo.Promo) {
	channels := b.config.getChannels()
	for _, ch := range channels {
		if _, err := b.sess.Channel(ch); err != nil {
			fmt.Println("Channel doesn't exist: ", ch)
			b.config.removeChannel(ch)
			continue
		}

		_, err := b.sess.ChannelMessageSendComplex(ch, m)
		if err != nil {
			fmt.Println("Post err: ", err)
		}
	}
}

func buildMessage(pr promo.Promo) *discordgo.MessageSend {
	meta := pr.GetMeta()

	embed := &discordgo.MessageEmbed{
		Title:       meta.Name,
		URL:         meta.Website,
		Description: meta.Description,
		Author: &discordgo.MessageEmbedAuthor{
			URL:     meta.Website,
			IconURL: meta.Icon,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: meta.Icon,
		},
		Fields: []*discordgo.MessageEmbedField{},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "#promotion",
		},
	}

	setPromoType(embed, pr)
	addWebsites(embed, *meta)
	addTags(embed, *meta)
	addMedia(embed, *meta)

	return &discordgo.MessageSend{
		Content: "` `",
		Embed: embed,
	}
}

func setPromoType(embed *discordgo.MessageEmbed, pr promo.Promo) {
	switch pr.GetMeta().Type {
	case "server":
		s := pr.(*promo.Server)
		embed.Color = 0x00d56a
		embed.Author.Name = "#Server"
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:   "IP",
				Value:  fmt.Sprintf("`%s`", s.IP),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Whitelist",
				Value:  promo.Or(s.Whitelist, "`Yes`", "`No`"),
				Inline: true,
			},
		)
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
}

func addWebsites(embed *discordgo.MessageEmbed, pr promo.Meta)  {
	if pr.Website != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Website",
			Value:  fmt.Sprintf("[%s](%s)", pr.Website, pr.Website),
			Inline: true,
		})
	}
	if pr.Discord != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Discord",
			Value:  fmt.Sprintf("[#Join](%s)", pr.Discord),
			Inline: true,
		})
	}
}

func addMedia(embed *discordgo.MessageEmbed, pr promo.Meta) {
	if pr.Media.Type == "image" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: pr.Media.URL,
		}
	}
}

func addTags(embed *discordgo.MessageEmbed, pr promo.Meta)  {
	if len(pr.Tags) > 0 {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text:  "#" + strings.Join(pr.Tags, " #"),
		}
	}
}