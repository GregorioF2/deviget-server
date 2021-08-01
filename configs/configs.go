package configs

import (
	"fmt"
	"os"
	"strconv"
)

func getEnvWithDefault(env, def string) string {
	res := os.Getenv(env)
	if res == "" {
		return def
	}
	return res
}

func getEnvWithDefaultInt(env string, def int64) int64 {
	resStr := os.Getenv(env)
	if resStr == "" {
		return def
	}
	resInt, err := strconv.ParseInt(resStr, 10, 64)
	if err != nil {
		fmt.Printf("WARNING: Env variable '%s' has not valid value", env)
		return def
	}
	return resInt
}

var SERVER_PORT string = getEnvWithDefault("SERVER_PORT", "3000")

var CACHE_MAX_TIME int64 = getEnvWithDefaultInt("CACHE_MAX_TIME", 60)
