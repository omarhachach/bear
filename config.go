package bear

// Config holds the configuration options for the bot.
type Config struct {
	DiscordToken  string     `json:"discord_token"`
	CommandPrefix string     `json:"command_prefix"`
	AutoDelete    bool       `json:"auto_delete"`
	AdminRoles    []string   `json:"admin_roles"`
	AdminUsers    []string   `json:"admin_users"`
	Log           *LogConfig `json:"log"`
}

// LogConfig holds the configuration options for the logger.
type LogConfig struct {
	Debug bool   `json:"debug"`
	File  string `json:"file"`
}
