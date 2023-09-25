package main

type Config struct {
	Keycloak *Keycloak
	Logger   *LoggerConfig
}

type Keycloak struct {
	Host        string
	Realm       string
	AuthUrl     string
	TokenUrl    string
	UserInfoUrl string
}

type LoggerConfig struct {
	Level LoggerLevel
}

type LoggerLevel string

const (
	LoggerLevelDebug  LoggerLevel = "debug"
	LoggerLevelInfo   LoggerLevel = "info"
	LoggerLevelSilent LoggerLevel = "silent"
)

func newAppConfig() *Config {
	kcc := Keycloak{
		Host:  getEnv("KEYCLOAK_HOST", "http://localhost:8080"),
		Realm: getEnv("KEYCLOAK_REALM", "master"),
	}

	kcc.AuthUrl = kcc.Host + "/auth/realms/" + kcc.Realm + "/protocol/openid-connect/auth"
	kcc.TokenUrl = kcc.Host + "/auth/realms/" + kcc.Realm + "/protocol/openid-connect/token"
	kcc.UserInfoUrl = kcc.Host + "/auth/realms/" + kcc.Realm + "/protocol/openid-connect/userinfo"

	return &Config{
		Keycloak: &kcc,
		Logger: &LoggerConfig{
			Level: LoggerLevel(getEnv("LOG_LEVEL", "debug")),
		},
	}
}
