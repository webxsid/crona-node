package kernel

import (
	"bufio"
	"crona/shared/protocol"
	sharedtypes "crona/shared/types"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"crona/tui/internal/logger"
)

type Info struct {
	PID        int    `json:"pid"`
	Port       int    `json:"port"`
	SocketPath string `json:"socketPath"`
	Token      string `json:"token"`
	ScratchDir string `json:"scratchDir"`
	Env        string `json:"env"`
}

func Ensure() (*Info, error) {
	home, _ := os.UserHomeDir()
	infoPath := filepath.Join(home, ".crona", "kernel.json")

	if info, err := readInfo(infoPath); err == nil {
		if isHealthy(info) {
			logger.Infof("Kernel already running at %s (pid %d)", info.SocketPath, info.PID)
			return info, nil
		}
	}

	logger.Info("Spawning kernel...")
	if err := launch(); err != nil {
		return nil, fmt.Errorf("launch kernel: %w", err)
	}

	for i := 0; i < 20; i++ {
		time.Sleep(250 * time.Millisecond)
		if info, err := readInfo(infoPath); err == nil {
			if isHealthy(info) {
				logger.Infof("Kernel ready at %s", info.SocketPath)
				return info, nil
			}
		}
	}

	return nil, fmt.Errorf("kernel failed to start within 5s")
}

func readInfo(path string) (*Info, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw struct {
		PID        int    `json:"pid"`
		Port       int    `json:"port"`
		SocketPath string `json:"socketPath"`
		Token      string `json:"token"`
		ScratchDir string `json:"scratchDir"`
		Env        string `json:"env"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	return &Info{
		PID:        raw.PID,
		Port:       raw.Port,
		SocketPath: raw.SocketPath,
		Token:      raw.Token,
		ScratchDir: raw.ScratchDir,
		Env:        raw.Env,
	}, nil
}

func isHealthy(info *Info) bool {
	if info.PID > 0 {
		if proc, err := os.FindProcess(info.PID); err != nil || proc == nil {
			return false
		}
	}
	conn, err := net.DialTimeout("unix", info.SocketPath, 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(2 * time.Second))

	reqBody, err := json.Marshal(protocol.Request{
		ID:     "healthcheck",
		Method: protocol.MethodHealthGet,
	})
	if err != nil {
		return false
	}
	if _, err := conn.Write(append(reqBody, '\n')); err != nil {
		return false
	}

	var resp protocol.Response
	if err := json.NewDecoder(bufio.NewReader(conn)).Decode(&resp); err != nil {
		return false
	}
	if resp.Error != nil {
		return false
	}
	var health sharedtypes.Health
	if err := json.Unmarshal(resp.Result, &health); err != nil {
		return false
	}
	return health.DB && health.Status == "ok"
}

func launch() error {
	cmd := exec.Command("crona-kernel")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	return cmd.Start()
}
