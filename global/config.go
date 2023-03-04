package global

type Configuration struct {
	ApiKey     string `mapstructure:"api_key"`
	CustomUrl  string `mapstructure:"custom_url"`
	Model      string `mapstructure:"model"`
	Proxy      string `mapstructure:"proxy"`
	ServerPort string `mapstructure:"server_port"`
	Timeout    int    `mapstructure:"timeout"`
	Version    string `mapstructure:"version"`

	Log Log `mapstructure:"log"`
}
