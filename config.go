package bear

// Config holds the configuration options for the bot.
type Config struct {
	Log           *LogConfig `json:"log"`
	DiscordToken  string     `json:"discord_token"`
	CommandPrefix string     `json:"command_prefix"`
	AutoDelete    bool       `json:"auto_delete"`
}

type LogConfig struct {
	Debug bool   `json:"debug"`
	File  string `json:"file"`
}
