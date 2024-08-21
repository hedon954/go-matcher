package config

// ServerConfig defines the server config.
type ServerConfig struct {
	AsynqRedis   *RedisOpt            `yaml:"asynq_redis"`
	NacosServers []*NacosServerConfig `yaml:"nacos_servers"`
}

type RedisOpt struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type NacosServerConfig struct {
	Addr        string `yaml:"addr"`
	Port        int    `yaml:"port"`
	ContextPath string `yaml:"context_path"`
	Schema      string `yaml:"schema"`
}
