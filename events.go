package bear

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func onMessageCreate(b *Bear) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(session *discordgo.Session, m *discordgo.MessageCreate) {
		// If sender is the bot.
		if m.Author.ID == session.State.User.ID {
			return
		}

		msg := m.ContentWithMentionsReplaced()
		caller := strings.Split(msg, " ")[0]

		cmd := b.Commands[caller]
		if cmd == nil {
			return
		}

		if b.Config.AutoDelete {
			_ = session.ChannelMessageDelete(m.ChannelID, m.ID)
		}

		go cmd.Handler(&Context{
			Session:   b.Session,
			Message:   m,
			ChannelID: m.ChannelID,
		})
	}
}
