package util

import "strings"

// IsCoreEnvVariable determines whether or not an environment variable key is a core Dokku variables.
// Core variables should _not_ be settable or viewable by clients.
func IsCoreEnvVariable(key string) bool {
	if key == "GIT_REV" || key == "PORT" || strings.HasPrefix(key, "DOKKU_") {
		return true
	}

	return false
}
