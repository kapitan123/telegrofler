package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	telegramTokenEnv  = "TELEGRAM_BOT_TOKEN"
	projectIdEnv      = "FIRESTORE_PROJECT_ID"
	gcloudAppCredsEnv = "GOOGLE_APPLICATION_CREDENTIALS"
)

var (
	TelegramToken = os.Getenv(telegramTokenEnv)
	ProjectId     = os.Getenv(projectIdEnv)
	ServerPort    = 9001 // AK TODO pass in env var
	GcloudCreds   = os.Getenv(gcloudAppCredsEnv)
	WorkersCount  = 1
)

func init() {
	if TelegramToken == "" {
		log.Panic("Telegram bot token is not set. Please set the environment variable ", telegramTokenEnv)
	}

	if ProjectId == "" {
		log.Panic("Firestore projectid is not set. Please set the environment variable ", projectIdEnv)
	}

	if GcloudCreds == "" {
		log.Info("gcloud creds not set. ADC default will be used. Variable ", gcloudAppCredsEnv)
	}
}
