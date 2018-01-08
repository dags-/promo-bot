package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strconv"
	"strings"
)

var numMatch = regexp.MustCompile("!interval ([0-9]+) .*?")

func (b *Bot) ready(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("Bot ready")
	s.UpdateStatus(0, "")
	b.StartLoop()
}

func (b *Bot) join(s *discordgo.Session, j *discordgo.GuildCreate) {
	fmt.Println("Joined guild", j.Name)
}

func (b *Bot) command(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.Author.ID != owner {
		return
	}

	mentions := false
	for _, m := range m.Mentions {
		if m.ID == s.State.User.ID {
			mentions = true
			break
		}
	}

	if !mentions {
		return
	}

	input := strings.ToLower(m.Content)
	if strings.Contains(input, "!add") && !b.config.hasChannel(m.ChannelID) {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		b.config.addChannel(m.ChannelID)
		fmt.Printf("Adding channel: %s\n", m.ChannelID)
		return
	}

	if strings.Contains(input, "!remove") && b.config.hasChannel(m.ChannelID) {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		b.config.removeChannel(m.ChannelID)
		fmt.Printf("Removing channel: %s\n", m.ChannelID)
		return
	}

	if strings.Contains(input, "!interval ") && b.config.hasChannel(m.ChannelID) {
		match := numMatch.FindStringSubmatch(input)
		if len(match) != 2 {
			fmt.Printf("Invalid time: '%s'\n", input)
			return
		}

		num, err := strconv.Atoi(match[1])
		if err != nil {
			fmt.Println(err)
		}

		s.ChannelMessageDelete(m.ChannelID, m.ID)
		b.config.setInterval(num)
		fmt.Printf("Setting interval: %v\n", num)
	}
}