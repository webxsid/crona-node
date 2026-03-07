import { useEffect, useState } from "react";
import type { KernelInfo } from "../interfaces/kernel.interface.js";
import type { ScratchPadMeta } from "../interfaces/scratchpad_meta.interface.js";
import type { KernelEvent } from "../interfaces/events.interface.js";

import { getScratchpads } from "../services/scratchpad-api.js";
import { eventBus } from "../state/event-bus.js";

export function useScratchpads(kernel: KernelInfo) {

  const [scratchpads, setScratchpads] = useState<ScratchPadMeta[]>([]);

  // initial load
  useEffect(() => {
    async function load() {
      try {
        const data = await getScratchpads(kernel);
        setScratchpads(data);
      } catch (err) {
        console.error("Failed to load scratchpads", err);
      }
    }

    load();
  }, [kernel.baseUrl]);

  // live updates
  useEffect(() => {

    const unsubscribe = eventBus.subscribe((event: KernelEvent) => {

      if (event.type === "scratchpad.created") {
        setScratchpads(prev => [...prev, event.payload]);
      }

      if (event.type === "scratchpad.updated") {
        setScratchpads(prev =>
          prev.map(s =>
            s.id === event.payload.id ? event.payload : s
          )
        );
      }

      if (event.type === "scratchpad.deleted") {
        setScratchpads(prev =>
          prev.filter(s => s.id !== event.payload.id)
        );
      }

    });

    return unsubscribe;

  }, []);

  return scratchpads;
}
