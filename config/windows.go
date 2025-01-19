//go:build windows

package config

func getNginxConfigDir(_ string) string {
	return "C:\\nginx\\conf\\conf.d"
}
