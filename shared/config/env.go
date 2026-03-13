package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const (
	EnvVarMode = "CRONA_ENV"
	ModeProd   = "Prod"
	ModeDev    = "Dev"
)

type AppEnv struct {
	Mode string
}

func Load() AppEnv {
	_ = loadDotEnvUp(".env")
	return Current()
}

func Current() AppEnv {
	mode := strings.TrimSpace(os.Getenv(EnvVarMode))
	if strings.EqualFold(mode, ModeDev) {
		return AppEnv{Mode: ModeDev}
	}
	return AppEnv{Mode: ModeProd}
}

func (e AppEnv) IsDev() bool {
	return e.Mode == ModeDev
}

func loadDotEnvUp(name string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	for {
		path := filepath.Join(wd, name)
		if _, err := os.Stat(path); err == nil {
			return loadDotEnvFile(path)
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return nil
		}
		wd = parent
	}
}

func loadDotEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		_ = os.Setenv(key, value)
	}
	return scanner.Err()
}
