package config

// App конфиг приложения
type App struct {
	Name           string `envconfig:"APP_NAME" default:"binance_bot"`
	Environment    string `envconfig:"APP_ENV" default:"prod"`
	Host           string `envconfig:"APP_HOST" default:"localhost"`
	Port           string `envconfig:"APP_PORT" default:"8081"`
	LogLevel       string `envconfig:"APP_LOG_LEVEL" default:"info"`
	PrettyLogs     bool   `envconfig:"APP_LOG_PRETTY" default:"false"`
	SwaggerFolder  string `envconfig:"SWAGGER_FOLDER" default:"swagger"`
	TZ             string `envconfig:"TZ" default:"Europe/Moscow"`
	GinMode        string `envconfig:"GIN_MODE" default:"release"`
	TelegramToken  string `envconfig:"TELEGRAM_TOKEN" default:""`
	TelegramUserID int    `envconfig:"TELEGRAM_USER_ID" default:""`
	TelegramHost   string `envconfig:"TELEGRAM_HOST" default:"api.telegram.org"`
	BatchSize      int    `envconfig:"BATCH_SIZE" default:"100"`
	EncryptedUid   string `envconfig:"ENCRYPTED_UID" default:""`
	TradeType      string `envconfig:"TRADE_TYPE" default:""`
}

func GetAppConfig(config *Config) App {
	return config.App
}
