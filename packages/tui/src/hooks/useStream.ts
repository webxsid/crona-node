import { useEffect, useState } from "react";
import type { KernelInfo } from "../interfaces/kernel.interface.js";
import type { Stream } from "../interfaces/stream.interface.js";
import { getStreams } from "../services/stream-api.js";
import type { KernelEvent } from "../interfaces/events.interface.js";
import { eventBus } from "../state/event-bus.js";

export function useStreams(
  kernel: KernelInfo,
  repoId?: string
) {
  const [streams, setStreams] = useState<Stream[]>([]);

  useEffect(() => {
    if (!repoId) {
      setStreams([]);
      return;
    }

    async function load() {
      try {
        const data = await getStreams(kernel, repoId);
        setStreams(data);
      } catch (err) {
        console.error("Failed to load streams", err);
      }
    }

    load();
  }, [kernel.baseUrl, repoId]);

  // Live updates from kernel
  useEffect(() => {
    const unsubscribe = eventBus.subscribe((event: KernelEvent) => {

      if (event.type === "stream.created") {
        if (event.payload.repoId !== repoId) return;

        setStreams((prev) => {
          // Avoid duplicates
          if (prev.some((s) => s.id === event.payload.id)) {
            return prev;
          }
          return [...prev, event.payload];
        });
      }

      if (event.type === "stream.updated") {
        setStreams((prev) =>
          prev.map((s) =>
            s.id === event.payload.id ? event.payload : s
          )
        );
      }

      if (event.type === "stream.deleted") {
        setStreams((prev) =>
          prev.filter((s) => s.id !== event.payload.id)
        );
      }

    });

    return unsubscribe;
  }, [repoId]);


  return streams;
}
