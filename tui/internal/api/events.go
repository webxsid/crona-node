package api

import (
	"bufio"
	"crona/shared/protocol"
	"encoding/json"
	"net"
	"strings"
	"time"
)

// Subscribe connects to the kernel event stream over the Unix socket.
// The goroutine reconnects on disconnect. Close the done channel to stop.
func Subscribe(socketPath string, done <-chan struct{}) <-chan KernelEvent {
	ch := make(chan KernelEvent, 32)

	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			default:
			}

			err := readStream(socketPath, ch, done)
			if err != nil {
				select {
				case <-done:
					return
				case <-time.After(2 * time.Second):
				}
			}
		}
	}()

	return ch
}

func readStream(socketPath string, ch chan<- KernelEvent, done <-chan struct{}) error {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	reqBody, err := json.Marshal(protocol.Request{
		ID:     "events-subscribe",
		Method: protocol.MethodEventsSubscribe,
	})
	if err != nil {
		return err
	}
	if _, err := conn.Write(append(reqBody, '\n')); err != nil {
		return err
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		select {
		case <-done:
			return nil
		default:
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var frame protocol.Event
		if err := json.Unmarshal([]byte(line), &frame); err != nil {
			continue
		}
		event := KernelEvent{Type: frame.Type, Payload: frame.Payload}

		select {
		case ch <- event:
		case <-done:
			return nil
		}
	}

	return scanner.Err()
}
