package config

const (
	// redis
	RedisAddr string = "localhost:6379"
	MaxExTime int    = 600

	// mysql
	DSN          string = "root@tcp(127.0.0.1:3306)/ums_dev"
	MaxIdleConns int    = 2000
	MaxOpenConns int    = 2000

	// TCP Server
	TCPServerAddr     string = "[127.0.0.1]:8888"
	TCPClientPoolSize int    = 2000

	// HTTP Server
	HTTPServerAddr   string = "127.0.0.1:8080"
	PprofAddr        string = "127.0.0.1:9999"
	StaticFilePath   string = "./static/images"
	DefaultImagePath string = "500.jpeg"
)
