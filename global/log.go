package global

type Log struct {
	Level      string `mapstructure:"level"`
	RootDir    string `mapstructure:"root_dir"`
	Filename   string `mapstructure:"filename"`
	Format     string `mapstructure:"format"`
	ShowLine   bool   `mapstructure:"show_line"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxSize    int    `mapstructure:"max_size"` // MB
	MaxAge     int    `mapstructure:"max_age"`  // day
	Compress   bool   `mapstructure:"compress"`
}
