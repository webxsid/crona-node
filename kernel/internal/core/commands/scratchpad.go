package commands

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"path"
	"regexp"
	"strings"
	"time"

	"crona/kernel/internal/core"

	"github.com/google/uuid"

	sharedtypes "crona/shared/types"
)

var (
	variablePattern        = regexp.MustCompile(`\[\[(\w+)\]\]`)
	scratchpadInvalidChars = regexp.MustCompile(`[<>:"\\|?*]`)
)

func ValidateScratchpadPath(inputPath string) error {
	if strings.TrimSpace(inputPath) == "" {
		return errors.New("path cannot be empty")
	}
	normalized := path.Clean(strings.TrimSpace(inputPath))
	if strings.HasPrefix(normalized, "/") {
		return errors.New("absolute paths are not allowed")
	}
	if strings.HasPrefix(normalized, "../") || strings.Contains(normalized, "/../") {
		return errors.New("path traversal is not allowed")
	}
	for _, segment := range strings.Split(normalized, "/") {
		if segment == "" {
			continue
		}
		if segment == "." || segment == ".." {
			return errors.New("invalid path segment")
		}
		if scratchpadInvalidChars.MatchString(segment) {
			return errors.New(`path contains invalid characters: < > : " \ | ? *`)
		}
	}
	return nil
}

func RegisterScratchpad(ctx context.Context, c *core.Context, meta sharedtypes.ScratchPadMeta) (string, error) {
	if strings.TrimSpace(meta.Path) == "" {
		return "", errors.New("path is required to register a scratchpad")
	}
	normalizedPath := strings.TrimSpace(meta.Path)
	if err := ValidateScratchpadPath(normalizedPath); err != nil {
		return "", err
	}
	processedPath, err := handleVariablesInPath(normalizedPath)
	if err != nil {
		return "", err
	}
	lastOpened := c.Now()
	if meta.LastOpenedAt != "" {
		lastOpened = meta.LastOpenedAt
	}
	id := meta.ID
	if id == "" {
		id = uuid.NewString()
	}
	if err := c.ScratchPads.Upsert(ctx, sharedtypes.ScratchPadMeta{
		ID:           id,
		Name:         meta.Name,
		Path:         processedPath,
		Pinned:       meta.Pinned,
		LastOpenedAt: lastOpened,
	}, c.UserID, c.DeviceID); err != nil {
		return "", err
	}
	return processedPath, nil
}

func GetScratchpad(ctx context.Context, c *core.Context, id string) (*sharedtypes.ScratchPadMeta, error) {
	return c.ScratchPads.GetByID(ctx, id, c.UserID, c.DeviceID)
}

func ListScratchpads(ctx context.Context, c *core.Context, pinnedOnly bool) ([]sharedtypes.ScratchPadMeta, error) {
	return c.ScratchPads.List(ctx, c.UserID, c.DeviceID, pinnedOnly)
}

func PinScratchpad(ctx context.Context, c *core.Context, id string, pinned bool) error {
	existing, err := c.ScratchPads.GetByID(ctx, id, c.UserID, c.DeviceID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("scratchpad not found")
	}
	existing.Pinned = pinned
	return c.ScratchPads.Upsert(ctx, *existing, c.UserID, c.DeviceID)
}

func RemoveScratchpad(ctx context.Context, c *core.Context, id string) error {
	return c.ScratchPads.RemoveByID(ctx, id, c.UserID, c.DeviceID)
}

func handleVariablesInPath(input string) (string, error) {
	allowed := map[string]bool{
		"date":      true,
		"time":      true,
		"datetime":  true,
		"timestamp": true,
		"random":    true,
	}
	found := variablePattern.FindAllStringSubmatch(input, -1)
	for _, match := range found {
		if len(match) < 2 || !allowed[match[1]] {
			return "", fmt.Errorf("invalid variable in path")
		}
	}
	now := time.Now()
	result := input
	for _, match := range found {
		var replacement string
		switch match[1] {
		case "date":
			replacement = now.Format("2006-01-02")
		case "time":
			replacement = now.Format("15-04-05")
		case "datetime":
			replacement = now.UTC().Format("2006-01-02_15-04-05")
		case "timestamp":
			replacement = fmt.Sprintf("%d", now.UnixMilli())
		case "random":
			replacement = randomString(8)
		}
		result = strings.Replace(result, match[0], replacement, 1)
	}
	return result, nil
}

func randomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteByte(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
