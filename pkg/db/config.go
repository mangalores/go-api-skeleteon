package db

type Config struct {
	Host         string `envconfig:"DB_HOST"`
	Driver       string `envconfig:"DB_DRIVER"`
	DatabaseName string `envconfig:"DB_NAME"`
	User         string `envconfig:"DB_USER"`
	Password     string `envconfig:"DB_PASSWORD"`
	Port         string `envconfig:"DB_PORT"`
	Logging      string `envconfig:"DB_LOGGING"`
}
