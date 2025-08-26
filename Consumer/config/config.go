package config

type EMQXConfig struct {
	Brokers                []string `yaml:"brokers"`
	Username               string   `yaml:"username"`
	Password               string   `yaml:"password"`
	ClientIDPrefix         string   `yaml:"client_id_prefix"`
	ShardSubscriptionGroup string   `yaml:"shard_subscription_group"`
	Qos                    int      `yaml:"qos"`
	Topics                 []string `yaml:"topics"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

func LoadConfig(path string) {
	// viper.AddConfigPath(path)
	// viper.SetConfigName("config")
	// viper.SetConfigType("yaml")

	// if err := viper.ReadInConfig(); err != nil {
	// 	return EMQXConfig{}, DatabaseConfig{}, LoggingConfig{}, err
	// }

}
