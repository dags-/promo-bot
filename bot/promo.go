package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dags-/promo-bot/promo"
	"github.com/dags-/promo-bot/util"
	"strings"
	"time"
)

func (b *Bot) StartLoop() {
	go func() {
		h := NewHistory(5)

		for {
			promos := b.api.GetPromoQueue()
			start := time.Now()

			fmt.Printf("Starting new promotions run, count=%v\n", len(promos))
			for _, pr := range promos {
				if h.Contains(pr.ID) {
					continue
				}

				m := buildMessage(pr)
				b.sendToAll(m, pr)
				h.Add(pr.ID)
				time.Sleep(b.config.getInterval())
			}

			remaining := b.config.getInterval() - time.Since(start)
			if remaining > 0 {
				time.Sleep(remaining)
			}
		}
	}()
}

func (b *Bot) sendToAll(m *discordgo.MessageSend, pr promo.Promotion) {
	channels := b.config.getChannels()
	for _, ch := range channels {
		if _, err := b.sess.Channel(ch); err != nil {
			fmt.Println("Err bot.sendAll.chan: Channel doesn't exist: ", ch)
			b.config.removeChannel(ch)
			continue
		}

		_, err := b.sess.ChannelMessageSendComplex(ch, m)
		if err != nil {
			fmt.Println("Err bot.sendAll.Send: ", err)
		}
	}
}

func buildMessage(pr promo.Promotion) *discordgo.MessageSend {
	embed := &discordgo.MessageEmbed{
		Title:       pr.Name,
		URL:         pr.Website,
		Description: pr.Description,
		Author: &discordgo.MessageEmbedAuthor{
			URL:     pr.Website,
			IconURL: pr.Icon,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: pr.Icon,
		},
		Fields: []*discordgo.MessageEmbedField{},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "#promotion",
		},
	}

	setPromoType(embed, pr)
	addWebsites(embed, pr)
	addTags(embed, pr)
	addMedia(embed, pr)

	return &discordgo.MessageSend{
		Embed: embed,
	}
}

func setPromoType(embed *discordgo.MessageEmbed, pr promo.Promotion) {
	switch pr.Type {
	case "server":
		embed.Color = 0x00d56a
		embed.Author.Name = "#Server"
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:   "IP",
				Value:  fmt.Sprintf("`%s`", *pr.IP),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Whitelist",
				Value:  utils.Or(pr.Whitelist != nil && *pr.Whitelist, "`Yes`", "`No`"),
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

func addWebsites(embed *discordgo.MessageEmbed, pr promo.Promotion)  {
	if pr.Website != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   utils.Or(pr.Type == "server", "Website", "Channel"),
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

func addMedia(embed *discordgo.MessageEmbed, pr promo.Promotion) {
	if pr.Image != "" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: pr.Image,
		}
	}
}

func addTags(embed *discordgo.MessageEmbed, pr promo.Promotion)  {
	if pr.Tags != "" {
		tags := strings.Split(pr.Tags, " ")
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text:  "#" + strings.Join(tags, " #"),
		}
	}
}