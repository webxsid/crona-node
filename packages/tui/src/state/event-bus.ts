import type { KernelEvent } from "../interfaces/events.interface.js";

type Listener = (event: KernelEvent) => void;

class EventBus {

  private listeners: Set<Listener> = new Set();

  emit(event: KernelEvent) {
    for (const listener of this.listeners) {
      listener(event);
    }
  }

  subscribe(listener: Listener) {
    this.listeners.add(listener);

    return () => {
      this.listeners.delete(listener);
    };
  }
}

export const eventBus = new EventBus();
