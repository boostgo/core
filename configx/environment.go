package configx

import (
	"strings"

	"github.com/boostgo/core/convert/format"
	"github.com/boostgo/core/fsx"
)

const (
	EnvLocal      = "local"
	EnvDevelop    = "dev"
	EnvProduction = "prod"
)

const (
	ExtensionJson = ".json"
	ExtensionYaml = ".yaml"
)

var (
	extension    = ExtensionYaml
	environments = make(map[string]struct{})
)

func init() {
	environments[EnvLocal] = struct{}{}
	environments[EnvDevelop] = struct{}{}
	environments[EnvProduction] = struct{}{}
}

func SetExtension(ext string) {
	extension = ext
}

func Env() string {
	// local environment case
	if Local() {
		return EnvLocal
	}

	// develop environment case
	if Develop() {
		return EnvDevelop
	}

	return EnvProduction
}

func Local() bool {
	if GetBool("LOCAL") {
		return true
	}

	return fsx.FileExist(configFileName(EnvLocal))
}

func Develop() bool {
	return GetBool("DEBUG")
}

func Production() bool {
	return !Local() && !Develop()
}

func Path(projectCode ...string) string {
	localConfigFileName := configFileName(EnvLocal, projectCode...)
	if fsx.FileExist(localConfigFileName) {
		return localConfigFileName
	}

	if Develop() {
		return configFileName(EnvDevelop, projectCode...)
	}

	return configFileName(EnvProduction, projectCode...)
}

func EnvPath() string {
	const (
		defaultFileName = ".env"
		localFileName   = "local.env"
	)

	if fsx.FileExist(localFileName) {
		return localFileName
	}

	return defaultFileName
}

func configFileName(env string, projectCode ...string) string {
	var suffix string
	if len(projectCode) > 0 {
		suffix = "." + format.Code(projectCode[0])
	}

	return strings.Join([]string{"config/", env, suffix, extension}, "")
}
