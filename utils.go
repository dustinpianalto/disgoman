package disgoman

/* helpers.go:
 * Utility functions to make my life easier.
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// GetDefaultStatusManager returns a default Status Manager
func GetDefaultStatusManager() StatusManager {
	return StatusManager{
		[]string{
			"Golang!",
			"DiscordGo!",
			"Disgoman!",
		}, "10s"}
}

// CheckPermissions checks the channel and guild permissions to see if the member has the needed permissions
func CheckPermissions(session *discordgo.Session, memberID string, channel discordgo.Channel, perms Permission) bool {
	if perms == 0 {
		return true // If no permissions are required then just return true
	}

	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.ID == memberID {
			if overwrite.Allow&int(perms) != 0 {
				return true // If the channel has an overwrite for the user then true
			} else if overwrite.Deny&int(perms) != 0 {
				return false // If there is an explicit deny then false
			}
		}
	}

	member, err := session.State.Member(channel.GuildID, memberID)
	if err != nil {
		fmt.Println(err)
		return false
	}
	for _, roleID := range member.Roles {
		role, err := session.State.Role(channel.GuildID, roleID)
		if err != nil {
			fmt.Println(err)
		}

		for _, overwrite := range channel.PermissionOverwrites {
			if overwrite.ID == roleID {
				if overwrite.Allow&int(perms) != 0 {
					return true // If the channel has an overwrite for the role then true
				} else if overwrite.Deny&int(perms) != 0 {
					return false // If there is an explicit deny then false
				}
			}
		}

		if role.Permissions&int(PermissionAdministrator) != 0 {
			return true // If they are an administrator then they automatically have all permissions
		}

		if role.Permissions&int(perms) != 0 {
			return true // The role has the required permissions
		}
	}

	return false // Default to false
}
