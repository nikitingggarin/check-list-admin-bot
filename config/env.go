package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadFromEnv() (*Config, error) {
	godotenv.Load()

	cfg := &Config{}
	cfg.Telegram.Token = os.Getenv("TELEGRAM_BOT_TOKEN")
	cfg.Database.URL = os.Getenv("SUPABASE_URL")
	cfg.Database.Key = os.Getenv("SUPABASE_KEY")
	cfg.Server.Port = os.Getenv("PORT")

	setDefaults(cfg)
	return cfg, validate(cfg)
}
