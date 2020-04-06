package disgoman

/* structs.go
 * Contains structs used in Disgoman
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

import "github.com/bwmarrin/discordgo"

// Create a CommandManager which will hold the info and structures required for handling command messages
type CommandManager struct {
	Prefixes         PrefixesFunc
	Owners           []string
	StatusManager    StatusManager
	OnErrorFunc      OnErrorFunc
	Commands         map[string]*Command
	IgnoreBots       bool
	CheckPermissions bool
}

// Create a StatusManager which will update the status of the bot
type StatusManager struct {
	Values   []string
	Interval string
}

// Create a Command with all of the info specific to this command
type Command struct {
	Name                string
	Aliases             []string
	Description         string
	OwnerOnly           bool
	Hidden              bool
	RequiredPermissions Permission
	Invoke              CommandInvokeFunc
}

// Structure containing all the context that a command needs to run
type Context struct {
	Session *discordgo.Session
	Channel *discordgo.Channel
	Message *discordgo.Message
	User    *discordgo.User
	Guild   *discordgo.Guild
	Member  *discordgo.Member
	Invoked string
}
