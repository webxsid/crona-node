package runtime

import (
	"os"
	"path/filepath"
)

const (
	dirPerm  = 0o700
	filePerm = 0o600
)

type Paths struct {
	BaseDir       string
	DBPath        string
	ScratchDir    string
	LogsDir       string
	InfoFile      string
	SocketPath    string
	CurrentLogDir string
}

func ResolvePaths() (Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Paths{}, err
	}

	base := filepath.Join(home, ".crona")
	logs := filepath.Join(base, "logs")

	return Paths{
		BaseDir:       base,
		DBPath:        filepath.Join(base, "crona.db"),
		ScratchDir:    filepath.Join(base, "scratch"),
		LogsDir:       logs,
		InfoFile:      filepath.Join(base, "kernel.json"),
		SocketPath:    filepath.Join(base, "kernel.sock"),
		CurrentLogDir: filepath.Join(logs, dateStamp()),
	}, nil
}

func EnsurePaths(paths Paths) error {
	for _, dir := range []string{
		paths.BaseDir,
		paths.ScratchDir,
		paths.LogsDir,
		paths.CurrentLogDir,
	} {
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return err
		}
	}
	return nil
}

func FilePerm() os.FileMode {
	return filePerm
}
