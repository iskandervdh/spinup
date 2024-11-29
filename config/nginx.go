package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var nginxConfigPath = "/etc/nginx/conf.d"

func restartNginx() {
	exec.Command("sudo", "systemctl", "restart", "nginx").Run()
}

func AddNginxConfig(name string, domain string, port int) error {
	config := fmt.Sprintf(`server {
	listen 80;

	server_name %s;

	location / {
		proxy_pass http://127.0.0.1:%d/;
		proxy_set_header Host $host;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Proto $scheme;
	}
}
`, domain, port)

	configPath := fmt.Sprintf("%s/%s.conf", nginxConfigPath, name)

	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file %s already exists", configPath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if config file exists: %v", err)
	}

	// Create the file
	err := exec.Command("sudo", "touch", configPath).Run()

	if err != nil {
		return err
	}

	// Write the config to the file
	createCommand := exec.Command("sudo", "tee", configPath)
	createCommand.Stdin = strings.NewReader(config)
	err = createCommand.Run()

	if err != nil {
		return err
	}

	restartNginx()

	return nil
}

func RemoveNginxConfig(name string) error {
	configPath := fmt.Sprintf("%s/%s.conf", nginxConfigPath, name)
	err := exec.Command("sudo", "rm", configPath).Run()

	if err != nil {
		return err
	}

	restartNginx()

	return nil
}
