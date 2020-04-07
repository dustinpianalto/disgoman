package disgoman

/* command-manager.go
 * The main command manager code
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kballard/go-shellquote"
	"log"
	"strings"
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

// OnMessage
// checks if the message has one of the specified prefixes
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

	// If we found our prefix then remove it and split the command into pieces
	cmd, err := shellquote.Split(strings.TrimPrefix(content, prefix))
	if err != nil {
		log.Fatal(err)
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

	channel, err := session.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Couldn't retrieve Channel.")
		return
	}

	if !CheckPermissions(session, *m.Member, *channel, command.RequiredPermissions) {
		embed := &discordgo.MessageEmbed{
			Title:       "Insufficient Permissions",
			Description: "You don't have the correct permissions to run this command.",
			Color:       0xFF0000,
		}
		if !command.Hidden {
			_, err := session.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				log.Println(err)
			}
		}
		return
	}

	me, err := session.GuildMember(m.GuildID, session.State.User.ID)
	if err != nil {
		log.Fatal(err)
		return
	}

	if !CheckPermissions(session, *me, *channel, command.RequiredPermissions) {
		embed := &discordgo.MessageEmbed{
			Title:       "Insufficient Permissions",
			Description: "I don't have the correct permissions to run this command.",
			Color:       0xFF0000,
		}
		if !command.Hidden {
			_, err := session.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				log.Println(err)
			}
		}
		return
	}

	if command.OwnerOnly && !c.IsOwner(m.Author.ID) {
		embed := &discordgo.MessageEmbed{
			Title:       "You can't run that command!",
			Description: "Sorry, only the bot owner(s) can run that command!",
			Color:       0xff0000,
		}

		if !command.Hidden {
			_, err := session.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				log.Println(err)
			}
		}
		return
	}

	guild, _ := session.Guild(m.GuildID)

	context := Context{
		Session: session,
		Channel: channel,
		Message: m.Message,
		User:    m.Author,
		Guild:   guild,
		Member:  m.Member,
		Invoked: invoked,
	}

	err = command.Invoke(context, cmd[1:])
	if err != nil && c.OnErrorFunc != nil {
		c.OnErrorFunc(context, cmd[0], err)
	}

}
