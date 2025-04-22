package config

// Config 应用配置结构
type Config struct {
	Server ServerConfig
	HBase  HBaseConfig
}

// ServerConfig Web服务器配置
type ServerConfig struct {
	Port string
}

// HBaseConfig HBase连接配置
type HBaseConfig struct {
	ZookeeperQuorum string
	Table           string
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		HBase: HBaseConfig{
			ZookeeperQuorum: "192.168.2.154:2181",
			Table:           "moviedata",
		},
	}
}
