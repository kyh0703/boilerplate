package configs

type DB struct {
	User         string   `mapstructure:"user"`
	Password     string   `mapstructure:"password"`
	SourceAddrs  []string `mapstructure:"sourceAddrs"`
	ReplicaAddrs []string `mapstructure:"replicaAddrs"`
	DBName       string   `mapstructure:"dbName"`
	FilePath     string   `mapstructure:"filePath"`
}

type Redis struct {
	MasterName    string   `mapstructure:"masterName"`
	SentinelAddrs []string `mapstructure:"sentinelAddrs"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers"`
}

type Sentry struct {
	Dsn string `mapstructure:"dsn"`
}
