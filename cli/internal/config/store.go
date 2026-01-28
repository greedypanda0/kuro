package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UserConfigPath returns the path to the user config file: ~/.kuro/config.
func UserConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".kuro", "config"), nil
}

func getFileConfig(flag int) (*os.File, error) {
	path, err := UserConfigPath()
	if err != nil {
		return nil, err
	}

	if flag&(os.O_WRONLY|os.O_RDWR|os.O_APPEND|os.O_CREATE) != 0 {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, err
		}
	}

	return os.OpenFile(path, flag, 0o644)
}

// GetUserConfig reads the user config file and returns the value for key.
// If the key does not exist, it returns os.ErrNotExist.
func GetUserConfig(key string) (string, error) {
	file, err := getFileConfig(os.O_RDONLY)
	if err != nil {
		if os.IsNotExist(err) {
			return "", os.ErrNotExist
		}
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		k := strings.TrimSpace(parts[0])
		if k == key {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", os.ErrNotExist
}

// SetUserConfig appends key/value to the user config file.
func SetUserConfig(key, value string) error {
	file, err := getFileConfig(os.O_CREATE | os.O_WRONLY | os.O_APPEND)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s=%s\n", key, value)
	return err
}
