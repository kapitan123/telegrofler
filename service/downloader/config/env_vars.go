package config

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	VidoFilesBucketUrl    string `env:"VIDEO_FILES_BUCKET_URL"`
	Port                  string `env:"PORT" envDefault:"8080"`
	VideoConvertedTopicId string `env:"VIDEO_CONVERTED_TOPIC_ID"`
	ProjectId             string `env:"PROJECT_ID"`
	DebguMode             bool   `env:"DEBUG_MODE"`
}

func GetEnvVars() Config {
	cfg := Config{}
	env.Parse(&cfg)

	if cfg.VidoFilesBucketUrl == "" {
		log.Panic("VIDEO_FILES_BUCKET_URL is not set")
	}

	if cfg.VideoConvertedTopicId == "" {
		log.Info("VIDEO_CONVERTED_TOPIC_ID is not set")
	}

	if cfg.ProjectId == "" {
		log.Info("PROJECT_ID is not set")
	}

	return cfg
}
