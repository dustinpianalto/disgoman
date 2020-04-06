package disgoman

/* structs.go
 * Contains structs used in Disgoman
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

import "github.com/bwmarrin/discordgo"

// A CommandManager which will hold the info and structures required for handling command messages
type CommandManager struct {
	Prefixes         PrefixesFunc
	Owners           []string
	StatusManager    StatusManager
	OnErrorFunc      OnErrorFunc
	Commands         map[string]*Command
	IgnoreBots       bool
	CheckPermissions bool
}

// A StatusManager which will update the status of the bot
type StatusManager struct {
	Values   []string
	Interval string
}

// A Command with all of the info specific to this command
type Command struct {
	Name                string
	Aliases             []string
	Description         string
	OwnerOnly           bool
	Hidden              bool
	RequiredPermissions Permission
	Invoke              CommandInvokeFunc
}

// Context that a command needs to run
type Context struct {
	// Discordgo Session Object
	Session *discordgo.Session
	// Channel the command was sent in
	Channel *discordgo.Channel
	// Original Message
	Message *discordgo.Message
	// Sending User
	User *discordgo.User
	// Guild the command was sent in
	Guild *discordgo.Guild
	// Sending Member
	Member *discordgo.Member
	// Name of the command as it was invoked (this is so you know what alias was used to call the command)
	Invoked string
}
