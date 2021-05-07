package util

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

// IsCoreEnvVariable determines whether or not an environment variable key is a core Dokku variables.
// Core variables should _not_ be settable or viewable by clients.
func IsCoreEnvVariable(key string) bool {
	if key == "GIT_REV" || key == "PORT" || key == "NO_VHOST" || strings.HasPrefix(key, "DOKKU_") {
		return true
	}

	return false
}

// Formats an environment map for use in a Dokku command
func FormatEnv(env map[string]string) []string {
	formatted_env := []string{}

	for key, value := range env {
		formatted_env = append(formatted_env, fmt.Sprintf(`%s=%s`, key,
			base64.StdEncoding.EncodeToString([]byte(value)),
		))
	}

	return formatted_env
}

func VerifyEnv(key string) bool {
	r, _ := regexp.Compile("^[a-zA-Z_][a-zA-Z0-9_]*$")
	return r.MatchString(key)
}
