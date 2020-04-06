package disgoman

/* structs.go
 * Contains structs used in Disgoman
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

import "github.com/bwmarrin/discordgo"

type CommandManager struct {
	Prefixes         PrefixesFunc
	Owners           []string
	StatusManager    StatusManager
	OnErrorFunc      OnErrorFunc
	Commands         map[string]*Command
	IgnoreBots       bool
	CheckPermissions bool
}

type StatusManager struct {
	Values   []string
	Interval string
}

type Command struct {
	Name                string
	Aliases             []string
	Description         string
	OwnerOnly           bool
	Hidden              bool
	RequiredPermissions Permission
	Invoke              CommandInvokeFunc
}

type Context struct {
	Session *discordgo.Session
	Channel *discordgo.Channel
	Message *discordgo.Message
	User    *discordgo.User
	Guild   *discordgo.Guild
	Member  *discordgo.Member
}
