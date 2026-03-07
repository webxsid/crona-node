import { useEffect } from "react";
import { eventBus } from "../state/event-bus.js";
import type { KernelEvent } from "../interfaces/events.interface.js";

export function useKernelEvent(handler: (event: KernelEvent) => void) {

  useEffect(() => {
    const unsubscribe = eventBus.subscribe(handler);
    return unsubscribe;
  }, [handler]);

}
