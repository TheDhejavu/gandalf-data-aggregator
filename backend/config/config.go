package config

type Config struct {
	Port string `env:"PORT" env-default:"8080"`

	Database struct {
		URL string `env-required:"true" env:"DATABASE_URL"`
	}

	Migrations struct {
		Path string `env-required:"true" env:"MIGRATIONS_PATH"`
	}

	ServerURL string `env-required:"true" env:"SERVER_URL"`

	JWTSecretKey string `env-required:"true" env:"JWT_SECRET_KEY"`

	Gandalf struct {
		PublicKey  string `env-required:"true" env:"GANDALF_APP_PUBLIC_KEY"`
		PrivateKey string `env-required:"true" env:"GANDALF_APP_PRIVATE_KEY"`
		SauronURL  string `env-required:"true" env:"GANDALF_SAURON_URL"`
	}

	Twitter struct {
		Key      string `env-required:"true" env:"TWITTER_KEY"`
		Secret   string `env-required:"true" env:"TWITTER_SECRET"`
		Callback string `env-required:"true" env:"TWITTER_CALLBACK"`
	}

	Redis struct {
		URL string `env-required:"true" env:"REDIS_URL"`
	}
}
