package bear

// Config holds the configuration options for the bot.
type Config struct {
	Debug         bool   `json:"debug"`
	DiscordToken  string `json:"discord_token"`
	CommandPrefix string `json:"command_prefix"`
	AutoDelete    bool   `json:"auto_delete"`
}
