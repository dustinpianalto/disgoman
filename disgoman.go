package disgoman

/* Package Disgoman:
 * Command Handler for DisgordGo.
 * Inspired by:
 *	- Anpan (https://github.com/MikeModder/anpan)
 *
 * Disgoman (c) 2020 Dusty.P/dustinpianalto
 */

func GetDefaultStatusManager() StatusManager {
	return StatusManager{
		[]string{
			"Golang!",
			"DiscordGo!",
			"Disgoman!",
		}, "10s"}
}
