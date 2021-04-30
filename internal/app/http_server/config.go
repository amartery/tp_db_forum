package http_server

// Config ...
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DataBaseURL string `toml:"database_url"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr:    ":5050",
		LogLevel:    "debug",
		DataBaseURL: "",
	}
}

// "host=database dbname=statserver_db user=statserver password=password sslmode=disable"
