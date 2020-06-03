package disgoman

/* command-manager.go
 * The main command manager code
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kballard/go-shellquote"
)

// AddCommand adds the Command at the address passed in to the Commands array on the CommandManager.
// This will return an error if the command's name or any of the aliases already exist.
func (c *CommandManager) AddCommand(command *Command) error {
	var aliases = []string{command.Name}
	if command.Aliases != nil {
		aliases = append(aliases, command.Aliases...)
	}
	for _, alias := range aliases {
		if _, ok := c.Commands[alias]; ok {
			return fmt.Errorf("An alias named %v already exists", alias)
		}
	}
	if len(aliases) > 0 {
		for _, alias := range aliases {
			c.Commands[alias] = command
		}
	}
	return nil
}

// RemoveCommand removes the command named from the Commands array
func (c *CommandManager) RemoveCommand(name string) error {
	deleted := false
	if _, ok := c.Commands[name]; ok {
		delete(c.Commands, name)
		deleted = true
	}
	if !deleted {
		return errors.New("command doesn't exist")
	}
	return nil
}

// IsOwner checks if the user ID passed in an owner of the bot
func (c *CommandManager) IsOwner(id string) bool {
	for _, o := range c.Owners {
		if o == id {
			return true
		}
	}
	return false
}

// OnMessage checks if the message has one of the specified prefixes
// and if the message contains one of the commands.
// It then processes the arguments to pass into the command,
// checks the permissions for the command, and
// runs the command function with the current Context.
func (c *CommandManager) OnMessage(session *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot && c.IgnoreBots {
		return // If the author is a bot and ignore bots is set then just exit
	}

	content := m.Content

	prefixes := c.Prefixes(m.GuildID)
	var prefix string
	var has bool
	for _, prefix = range prefixes {
		if strings.HasPrefix(content, prefix) {
			has = true
			break
		}
	}
	if !has {
		return // If we didn't find a valid prefix then exit
	}

	channel, err := session.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Couldn't retrieve Channel.")
		return
	}

	guild, _ := session.Guild(m.GuildID)

	var cmd []string
	// If we found our prefix then remove it and split the command into pieces
	cmd, err = shellquote.Split(strings.TrimPrefix(content, prefix))
	if err != nil {
		fmt.Println(err.Error())
		if strings.Contains(err.Error(), "Unterminated") {
			cmd = strings.Split(strings.TrimPrefix(content, prefix), " ")
		} else {
			ctx := Context{
				Session:      session,
				Channel:      channel,
				Message:      m.Message,
				User:         m.Author,
				Guild:        guild,
				Member:       m.Member,
				Invoked:      "",
				ErrorChannel: c.ErrorChannel,
			}
			c.ErrorChannel <- CommandError{
				Context: ctx,
				Message: "",
				Error:   err,
			}
			return
		}
	}

	if len(cmd) < 1 {
		return
	}

	var command *Command
	invoked := cmd[0]
	if cmnd, ok := c.Commands[invoked]; ok {
		command = cmnd
	} else {
		fmt.Println("Command Not Found")
		return
	}

	if command.SanitizeEveryone {
		for i := 1; i < len(cmd); i++ {
			cmd[i] = strings.ReplaceAll(cmd[i], "@everyone", "@\ufff0everyone")
			cmd[i] = strings.ReplaceAll(cmd[i], "@here", "@\ufff0here")
		}
	}

	if !CheckPermissions(session, m.Author.ID, *channel, command.RequiredPermissions) {
		ctx := Context{
			Session:      session,
			Channel:      channel,
			Message:      m.Message,
			User:         m.Author,
			Guild:        guild,
			Member:       m.Member,
			Invoked:      cmd[0],
			ErrorChannel: c.ErrorChannel,
		}
		c.ErrorChannel <- CommandError{
			Context: ctx,
			Message: "You don't have the correct permissions to run this command.",
			Error:   errors.New("insufficient permissions"),
		}
		return
	}

	if !CheckPermissions(session, session.State.User.ID, *channel, command.RequiredPermissions) {
		ctx := Context{
			Session:      session,
			Channel:      channel,
			Message:      m.Message,
			User:         m.Author,
			Guild:        guild,
			Member:       m.Member,
			Invoked:      cmd[0],
			ErrorChannel: c.ErrorChannel,
		}
		c.ErrorChannel <- CommandError{
			Context: ctx,
			Message: "I don't have the correct permissions to run this command.",
			Error:   errors.New("insufficient permissions"),
		}
		return

	}

	if command.OwnerOnly && !c.IsOwner(m.Author.ID) {
		ctx := Context{
			Session:      session,
			Channel:      channel,
			Message:      m.Message,
			User:         m.Author,
			Guild:        guild,
			Member:       m.Member,
			Invoked:      cmd[0],
			ErrorChannel: c.ErrorChannel,
		}
		c.ErrorChannel <- CommandError{
			Context: ctx,
			Message: "Sorry, only the bot owner(s) can run that command!",
			Error:   errors.New("insufficient permissions"),
		}
		return

	}

	context := Context{
		Session:      session,
		Channel:      channel,
		Message:      m.Message,
		User:         m.Author,
		Guild:        guild,
		Member:       m.Member,
		Invoked:      invoked,
		ErrorChannel: c.ErrorChannel,
	}

	go command.Invoke(context, cmd[1:])
}
