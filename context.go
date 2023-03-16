package bear

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Context holds useful functions for a handler.
type Context struct {
	ChannelID string
	Session   *discordgo.Session
	Message   *discordgo.MessageCreate
}

// This is a collection of the standard colors used for messages.
const (
	SuccessColor int = 5025616
	ErrorColor   int = 16007990
	InfoColor    int = 2201331
)

// SendEmbed sends a custom embed message.
func (ctx Context) SendEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSendEmbed(ctx.ChannelID, embed)
}

// SendMessage sends an embed message with a message and color.
func (ctx Context) SendMessage(color int, title, format string, a ...interface{}) (*discordgo.Message, error) {
	messageEmbed := &discordgo.MessageEmbed{
		Color: color,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  title,
				Value: fmt.Sprintf(format, a...),
			},
		},
	}

	return ctx.SendEmbed(messageEmbed)
}

// SendErrorMessage is a quick way to send an error.
func (ctx Context) SendErrorMessage(format string, a ...interface{}) (*discordgo.Message, error) {
	return ctx.SendMessage(ErrorColor, "Error", format, a...)
}

// SendSuccessMessage is a shortcut for sending a message with the success colors.
func (ctx Context) SendSuccessMessage(format string, a ...interface{}) (*discordgo.Message, error) {
	return ctx.SendMessage(SuccessColor, "Success", format, a...)
}

// SendInfoMessage is a shortcut for sending a message with the info colors.
func (ctx Context) SendInfoMessage(format string, a ...interface{}) (*discordgo.Message, error) {
	return ctx.SendMessage(InfoColor, "Info", format, a...)
}
