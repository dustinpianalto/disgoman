package disgoman

/* helpers.go:
 * Utility functions to make my life easier.
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sort"
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

// HasHigherRole checks if the caller has a higher role than the target in the given guild
func HasHigherRole(session *discordgo.Session, guildID string, callerID string, targetID string) (bool, error) {
	guild, _ := session.Guild(guildID)
	if targetID == guild.OwnerID {
		return false, nil // If the target is the guild owner then no one is higher than them
	}
	if callerID == guild.OwnerID {
		return true, nil // If the caller is the guild owner then they are automatically higher than everyone
	}
	caller, err := session.GuildMember(guildID, callerID)
	if err != nil {
		return false, err
	}
	target, err := session.GuildMember(guildID, targetID)
	if err != nil {
		return false, err
	}
	var callerRoles []*discordgo.Role
	for _, roleID := range caller.Roles {
		role, _ := session.State.Role(guildID, roleID)
		callerRoles = append(callerRoles, role)
	}
	sort.Slice(callerRoles, func(i, j int) bool { return callerRoles[i].Position > callerRoles[j].Position })
	var targetRoles []*discordgo.Role
	for _, roleID := range target.Roles {
		role, _ := session.State.Role(guildID, roleID)
		targetRoles = append(targetRoles, role)
	}
	sort.Slice(targetRoles, func(i, j int) bool { return targetRoles[i].Position > targetRoles[j].Position })
	if len(callerRoles) == 0 {
		return false, nil // If the caller doesn't have any roles then they can't be higher than the target
	}
	if len(callerRoles) > 0 && len(targetRoles) == 0 {
		return true, nil // If the caller has a role and the target doesn't then that role is higher
	}
	if callerRoles[0].Position > targetRoles[0].Position {
		return true, nil
	}
	return false, nil
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
