package disgoman

import (
	"io"

	"github.com/dustinpianalto/discordgo"
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

// SendEmbed sends an embed to originating channel
func (c *Context) SendEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSendEmbed(c.Channel.ID, embed)
}

// SendFile sends a file to originating channel
func (c *Context) SendFile(filename string, file io.Reader) (*discordgo.Message, error) {
	return c.Session.ChannelFileSend(c.Channel.ID, filename, file)
}

// TODO Combine these to all use ChannelMessageSendComplex

// SendError makes a CommandError and sends it to the ErrorChannel. This includes the current context in the error.
// Will block if the channel buffer is full. It is up to the client to implement a channel for the errors as well as
// a function to handle the errors from said channel. If the ErrorChannel is nil then this does nothing.
func (c *Context) SendError(message string, err error) {
	if c.CommandManager.ErrorChannel != nil {
		c.CommandManager.ErrorChannel <- CommandError{
			Context: *c,
			Message: message,
			Error:   err,
		}
	}
}
