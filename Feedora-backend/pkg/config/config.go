package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server        ServerConfig `yaml:"server"`
	MySQL         MySQLConfig  `yaml:"mysql"`
	JWT           JWTConfig    `yaml:"jwt"`
	OSS           OSSConfig    `yaml:"oss"`
	Kafka         KafkaConfig  `yaml:"kafka"`
	CORS          CORSConfig   `yaml:"cors"`
	Redis         RedisConfig  `yaml:"redis"`
	Elasticsearch ESConfig     `yaml:"elasticsearch"`
	Worker        WorkerConfig `yaml:"worker"`
}

type ServerConfig struct {
	Name string `yaml:"name"`
	Env  string `yaml:"env"`
	Port int    `yaml:"port"`
}

type MySQLConfig struct {
	DSN                    string `yaml:"dsn"`
	MaxOpenConns           int    `yaml:"maxOpenConns"`
	MaxIdleConns           int    `yaml:"maxIdleConns"`
	ConnMaxLifetimeSeconds int    `yaml:"connMaxLifetimeSeconds"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expireHours"`
}

type OSSConfig struct {
	Type          string `yaml:"type"`
	BasePath      string `yaml:"basePath"`
	PublicBaseUrl string `yaml:"publicBaseUrl"`
}

type KafkaConfig struct {
	Enabled       bool     `yaml:"enabled"`
	Brokers       []string `yaml:"brokers"`
	TopicPrefix   string   `yaml:"topicPrefix"`
	ConsumerGroup string   `yaml:"consumerGroup"`
}

type CORSConfig struct {
	AllowOrigins []string `yaml:"allowOrigins"`
}

// RedisConfig Redis 连接配置。
type RedisConfig struct {
	Enabled      bool   `yaml:"enabled"`
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"poolSize"`
	MinIdleConns int    `yaml:"minIdleConns"`
}

// ESConfig Elasticsearch 连接配置。
type ESConfig struct {
	Enabled     bool     `yaml:"enabled"`
	Addresses   []string `yaml:"addresses"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	IndexPrefix string   `yaml:"indexPrefix"`
}

// WorkerConfig Worker 运行配置。
type WorkerConfig struct {
	Enabled         bool `yaml:"enabled"`
	BatchSize       int  `yaml:"batchSize"`
	OutboxIntervalS int  `yaml:"outboxIntervalSeconds"`
	MaxRetry        int  `yaml:"maxRetry"`
}

// Load 从指定的 yaml 文件路径读取配置。
func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
