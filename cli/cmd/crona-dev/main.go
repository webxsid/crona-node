package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"crona/shared/config"
	"crona/shared/protocol"
)

func main() {
	_ = config.Load()

	if len(os.Args) < 2 {
		fail("usage: crona-dev <seed|clear>")
	}

	method := ""
	switch os.Args[1] {
	case "seed":
		method = protocol.MethodKernelSeedDev
	case "clear":
		method = protocol.MethodKernelClearDev
	default:
		fail("unknown command: " + os.Args[1])
	}

	socketPath, err := readSocketPath()
	if err != nil {
		fail(err.Error())
	}

	if err := call(socketPath, method); err != nil {
		fail(err.Error())
	}

	fmt.Printf("%s ok\n", os.Args[1])
}

func readSocketPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	body, err := os.ReadFile(filepath.Join(home, ".crona", "kernel.json"))
	if err != nil {
		return "", err
	}
	var info struct {
		SocketPath string `json:"socketPath"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return "", err
	}
	if info.SocketPath == "" {
		return "", fmt.Errorf("kernel socket path not found")
	}
	return info.SocketPath, nil
}

func call(socketPath string, method string) error {
	conn, err := net.DialTimeout("unix", socketPath, 5*time.Second)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))

	body, err := json.Marshal(protocol.Request{
		ID:     "crona-dev",
		Method: method,
	})
	if err != nil {
		return err
	}
	if _, err := conn.Write(append(body, '\n')); err != nil {
		return err
	}

	var resp protocol.Response
	if err := json.NewDecoder(bufio.NewReader(conn)).Decode(&resp); err != nil {
		return err
	}
	if resp.Error != nil {
		return fmt.Errorf("%s: %s", resp.Error.Code, resp.Error.Message)
	}
	return nil
}

func fail(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
