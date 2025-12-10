package config

type Config struct {
	Telegram struct {
		Token string `validate:"required"`
	}
	Database struct {
		URL string `validate:"required"`
		Key string `validate:"required"`
	}
	Server struct {
		Port string `default:":8080"`
	}
}
