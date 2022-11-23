package util

import (
	"os"
)

func Addr() string {
	if port, hasPort := os.LookupEnv("PORT"); hasPort {
		return ":" + port
	}
	return ":8080"
}
