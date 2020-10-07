package apiserver

// Параметры конфигурации
type Config struct {
	BindAddr	string `toml:"bindAddr"`
	DatabaseURL string `toml:"database_url"`
}

// Инициализация конфига с дефолтным значением
func NewConfig() *Config{
	return &Config{
		BindAddr: ":8080",
		DatabaseURL: "mongodb://127.0.0.1:27017",
	}
}