package controller

type Config struct {
	LogLevel              string `yaml:"log_level"`
	SecretID              string `yaml:"secret_id"`
	VersionID             string `yaml:"version_id"`
	Region                string
	Message               string
	MessageForSystemUser  string
	Slack                 SlackConfig
	InitialPasswordLength int `yaml:"initial_password_length"`
}

type SlackConfig struct {
	ChannelIDForSystemUser string `yaml:"channel_id_for_system_user"`
}

type Param struct {
	ConfigFilePath        string
	LogLevel              string
	SlackBotAccessToken   string
	MessageTemplateString string
	InitalPasswordLength  int
	// user filter
	// user mapping
	// notifier (slack, SES, etc)
	// integration (datadog, sentry, etc)
	DryRun   bool
	UserName string
}
