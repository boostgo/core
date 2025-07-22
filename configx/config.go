// Package configx helps to manipulate with app configs (yaml, json, env)
// Features:
// - Read yaml or json config file and load values to structure
// - Get environment variable as string, bool or int
// - Environment management - local, dev, prod. Getting config file path by current environment
// - Configuration samples. Common config structures like Server, Swagger, SQL, Redis, etc...
package configx

import (
	"os"
	"strconv"

	"github.com/boostgo/core/defaults"
	"github.com/ilyakaznacheev/cleanenv"
)

// Read export read config file to provided export object.
//
// Provided paths can contain as json/yaml file and also .env file
func Read(export any, path ...string) error {
	if len(path) == 0 {
		if err := cleanenv.ReadEnv(export); err != nil {
			return err
		}

		return defaults.Set(export)
	}

	for _, p := range path {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			continue
		}

		if err := cleanenv.ReadConfig(p, export); err != nil {
			return err
		}
	}

	return defaults.Set(export)
}

// MustRead calls Read function and if catch error throws panic
func MustRead(export any, path ...string) {
	if err := Read(export, path...); err != nil {
		panic(err)
	}
}

// GetString read environment variable and convert to string
func GetString(key string) string {
	return os.Getenv(key)
}

// GetBool read environment variable and convert to bool
func GetBool(key string) bool {
	switch os.Getenv(key) {
	case "true", "TRUE":
		return true
	default:
		return false
	}
}

// GetInt read environment variable and convert to int
func GetInt(key string) int {
	result, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0
	}

	return result
}
