package config

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	baseDirProd = ".crona"
	baseDirDev  = ".crona-dev"
)

func IsDevMode(mode string) bool {
	return strings.EqualFold(strings.TrimSpace(mode), ModeDev)
}

func RuntimeBaseDirNameForMode(mode string) string {
	if IsDevMode(mode) {
		return baseDirDev
	}
	return baseDirProd
}

func RuntimeBaseDirName() string {
	return RuntimeBaseDirNameForMode(Load().Mode)
}

func RuntimeBaseDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, RuntimeBaseDirName()), nil
}

func BinarySuffixForMode(mode string) string {
	if IsDevMode(mode) {
		return "-dev"
	}
	return ""
}

func CLIBinaryNameForMode(mode string) string {
	return "crona" + BinarySuffixForMode(mode)
}

func KernelBinaryNameForMode(mode string) string {
	return "crona-kernel" + BinarySuffixForMode(mode)
}

func TUIBinaryNameForMode(mode string) string {
	return "crona-tui" + BinarySuffixForMode(mode)
}

func CLIBinaryName() string {
	return CLIBinaryNameForMode(Load().Mode)
}

func KernelBinaryName() string {
	return KernelBinaryNameForMode(Load().Mode)
}

func TUIBinaryName() string {
	return TUIBinaryNameForMode(Load().Mode)
}
