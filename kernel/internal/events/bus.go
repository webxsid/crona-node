package events

import (
	"sync"

	sharedtypes "crona/shared/types"
)

type Listener func(sharedtypes.KernelEvent)

type Bus struct {
	mu        sync.RWMutex
	listeners map[int]Listener
	nextID    int
}

func NewBus() *Bus {
	return &Bus{
		listeners: make(map[int]Listener),
	}
}

func (b *Bus) Emit(event sharedtypes.KernelEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, listener := range b.listeners {
		listener(event)
	}
}

func (b *Bus) Subscribe(listener Listener) func() {
	b.mu.Lock()
	defer b.mu.Unlock()

	id := b.nextID
	b.nextID++
	b.listeners[id] = listener

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		delete(b.listeners, id)
	}
}
