package disgoman

import (
	"github.com/bwmarrin/discordgo"
	"io"
)

/* context.go:
 * Utility functions for command context
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

// Send message to originating channel
func (c *Context) Send(message string) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSend(c.Channel.ID, message)
}

// Send embed to originating channel
func (c *Context) SendEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSendEmbed(c.Channel.ID, embed)
}

// Send embed to originating channel
func (c *Context) SendFile(filename string, file io.Reader) (*discordgo.Message, error) {
	return c.Session.ChannelFileSend(c.Channel.ID, filename, file)
}

// TODO Combine these to all use ChannelMessageSendComplex
