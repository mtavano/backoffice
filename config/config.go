package config

type Config struct {
	Environment string `envconfig:"environment" required:"true"`
	SecretKey   string `envconfig:"environment" required:"true"`

	// Server config
	ApiServerPort string `envconfig:"port" required:"true"`

	// Database related
	DBDriver        string `envconfig:"db_driver" required:"true"`
	DBDsn           string `envconfig:"db_dsn" required:"true"`
	DBMigrationsDir string `envconfig:"db_migrations_dir" required:"true"`
}
