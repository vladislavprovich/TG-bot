package logger

type ConfigLogger struct {
	Level string `envconfig:"LOG_LEVEL" default:"info"`
}
