//go:build linux || darwin

package config

import (
	"path"
)

func getNginxConfigDir(configDir string) string {
	return path.Join(configDir, "nginx")
}
