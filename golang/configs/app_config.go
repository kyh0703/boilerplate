package configs

type Log struct {
	Level       string `mapstructure:"level"`
	HistoryType string `mapstructure:"historyType"`
}

type Server struct {
	Profile string `mapstructure:"profile"`
	Port    string `mapstructure:"port"`
}
