package disgoman

/* Package Disgoman:
 * Command Handler for DisgordGo.
 * Inspired by:
 *	- Anpan (https://github.com/MikeModder/anpan)
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

func GetCommandManager(prefixes PrefixesFunc, owners []string, ignoreBots, checkPerms bool) CommandManager {
	return CommandManager{
		Prefixes:         prefixes,
		Owners:           owners,
		StatusManager:    GetDefaultStatusManager(),
		Commands:         make(map[string]*Command),
		Aliases:          make(map[string]string),
		IgnoreBots:       ignoreBots,
		CheckPermissions: checkPerms,
	}
}

func GetStatusManager(values []string, interval string) StatusManager {
	return StatusManager{
		Values:   values,
		Interval: interval,
	}
}

func GetDefaultStatusManager() StatusManager {
	return GetStatusManager(
		[]string{
			"Golang!",
			"DiscordGo!",
			"Disgoman!",
		}, "10s")
}
