import { useEffect } from "react";
import { EventSource } from "eventsource";
import type { KernelInfo } from "../interfaces/kernel.interface.js";
import type { KernelEvent } from "../interfaces/events.interface.js";
import { eventBus } from "../state/event-bus.js";

export function useKernelEvents(kernel: KernelInfo) {

  useEffect(() => {
    const url = `${kernel.baseUrl}/events`;

    const es = new EventSource(url, {
      fetch: (url, init) => {
        init.headers = {
          ...init.headers,
          Authorization: `Bearer ${kernel.token}`
        };
        return fetch(url, init);
      }
    });

    es.onmessage = (e) => {
      try {
        const event: KernelEvent = JSON.parse(e.data);
        eventBus.emit(event);
      } catch (err) {
        console.error("Invalid kernel event:", err);
      }
    };

    es.onerror = (err) => {
      console.error("Kernel SSE connection error", err);
    };

    return () => {
      es.close();
    };

  }, [kernel.baseUrl, kernel.token]);
}
