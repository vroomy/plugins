package plugins

import (
	"path"
	"strings"
)

// ParseKey returns stripped gitUrl and plugin alias
func ParseKey(key string) (newKey, alias string) {
	spl := strings.Split(key, " as ")
	// Set key as the first part of the split
	newKey = spl[0]
	// Check to see if an alias was provided
	if len(spl) > 1 {
		// Alias was provided, set the alias value
		alias = spl[1]
	} else {
		_, name := path.Split(newKey)
		alias = strings.Split(name, "-")[0]
		alias = strings.Split(alias, "@")[0]
		alias = strings.Split(alias, "#")[0]
	}
	return
}
