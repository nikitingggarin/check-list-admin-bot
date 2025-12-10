package config

import "errors"

func validate(cfg *Config) error {
	if cfg.Telegram.Token == "" {
		return errors.New("TELEGRAM_BOT_TOKEN is required")
	}
	if cfg.Database.URL == "" {
		return errors.New("SUPABASE_URL is required")
	}
	if cfg.Database.Key == "" {
		return errors.New("SUPABASE_KEY is required")
	}
	return nil
}
