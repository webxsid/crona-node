package scratchfile

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func Create(root string, userPath string, title string) (string, error) {
	fullPath, err := resolve(root, userPath)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o700); err != nil {
		return "", err
	}
	f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fullPath, nil
		}
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()
	if safeTitle := strings.TrimSpace(title); safeTitle != "" {
		if _, err := f.WriteString("# " + safeTitle + "\n\n"); err != nil {
			return "", err
		}
	}
	return fullPath, nil
}

func Read(root string, userPath string) (string, error) {
	fullPath, err := resolve(root, userPath)
	if err != nil {
		return "", err
	}
	body, err := os.ReadFile(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", errors.New("scratchpad file not found")
		}
		return "", err
	}
	return string(body), nil
}

func Delete(root string, userPath string) error {
	fullPath, err := resolve(root, userPath)
	if err != nil {
		return err
	}
	err = os.Remove(fullPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func ClearAll(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return os.MkdirAll(root, 0o700)
		}
		return err
	}
	for _, entry := range entries {
		if err := os.RemoveAll(filepath.Join(root, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func resolve(root string, userPath string) (string, error) {
	if userPath == "" {
		return "", errors.New("scratchpad path is required")
	}
	normalized := userPath
	if !strings.HasSuffix(normalized, ".md") {
		normalized += ".md"
	}
	fullPath := filepath.Clean(filepath.Join(root, normalized))
	rootPath := filepath.Clean(root)
	if fullPath != rootPath && !strings.HasPrefix(fullPath, rootPath+string(os.PathSeparator)) {
		return "", errors.New("invalid scratchpad path")
	}
	return fullPath, nil
}
