package configs

type App struct {
	Server  Server `mapstructure:"server"`
	Version string `mapstructure:"version"`
	Log     Log    `mapstructure:"log"`
}

type Infra struct {
	DB     DB     `mapstructure:"db"`
	Redis  Redis  `mapstructure:"redis"`
	Kafka  Kafka  `mapstructure:"kafka"`
	Sentry Sentry `mapstructure:"sentry"`
}

type Auth struct {
	Google Google `mapstructure:"google"`
	Kakao  Kakao  `mapstructure:"kakao"`
	Github Github `mapstructure:"github"`
}

type Config struct {
	App    App    `mapstructure:"app"`
	Infra  Infra  `mapstructure:"infra"`
	Auth   Auth   `mapstructure:"auth"`
	Server Server `mapstructure:"server"`
}
