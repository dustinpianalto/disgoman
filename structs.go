package disgoman

/* structs.go
 * Contains structs used in Disgoman
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

import "github.com/dustinpianalto/discordgo"

// CommandManager holds the info and structures required for handling command messages
type CommandManager struct {
	// Function which returns a list of strings which are to be used as prefixes
	Prefixes PrefixesFunc
	// Array of string IDs of the owners of the bot (for use with commands that are owner only)
	Owners []string
	// Status Manager which will handle updating the status of the bot
	StatusManager StatusManager
	// Channel to send errors to
	ErrorChannel chan CommandError
	// Map of the command names to the pointer of the command to call
	Commands map[string]*Command
	// Should we ignore bots when processing commands
	IgnoreBots bool
	// Do we need to check permissions for commands (not used currently)
	CheckPermissions bool
}

// StatusManager updates the status of the bot
type StatusManager struct {
	// Values that will be used for the status
	Values []string
	// How often is the status updated
	Interval string
}

// Command contains all of the info specific to this command
type Command struct {
	// The name of the command
	Name string
	// Aliases the command can be called with
	Aliases []string
	// Simple description for use with the help command
	Description string
	// Is the command owner only
	OwnerOnly bool
	// Is the command hidden (should it show up in help)
	Hidden bool
	// Permissions that are required for the command to function
	RequiredPermissions Permission
	// Should the message content be sanitized before invoking the command
	SanitizeEveryone bool
	// Function to invoke when command is called
	Invoke CommandInvokeFunc
}

// CommandError contains all the information needed to process an error in a command
type CommandError struct {
	// The Context the command was run in
	Context Context
	// Error Message
	Message string
	// The Error object
	Error error
}

// Context contains all the context that a command needs to run
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
	// Command Manager
	CommandManager *CommandManager
}
