package bear

// Config holds the configuration options for the bot.
type Config struct {
	DiscordToken string   `json:"discord_token"`
	AutoDelete   bool     `json:"auto_delete"`
	AdminRoles   []string `json:"admin_roles"`
	AdminUsers   []string `json:"admin_users"`
}
