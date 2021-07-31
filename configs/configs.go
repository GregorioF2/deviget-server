package configs

import (
	"os"
)

func getEnvWithDefault(env, def string) string {
	res := os.Getenv(env)
	if res == "" {
		return def
	}
	return res
}

var SERVER_PORT string = getEnvWithDefault("SERVER_PORT", "3000")
