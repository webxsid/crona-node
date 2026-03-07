import { useEffect, useState } from "react";
import type { KernelInfo } from "../interfaces/kernel.interface.js";
import type { Issue } from "../interfaces/issue.interface.js";
import { getIssues } from "../services/issue-api.js";
import type { KernelEvent } from "../interfaces/events.interface.js";
import { eventBus } from "../state/event-bus.js";

export function useIssues(
  kernel: KernelInfo,
  streamId?: string
) {
  const [issues, setIssues] = useState<Issue[]>([]);

  useEffect(() => {
    if (!streamId) {
      setIssues([]);
      return;
    }

    async function load() {
      try {
        const data = await getIssues(kernel, streamId);
        setIssues(data);
      } catch (err) {
        console.error("Failed to load issues", err);
      }
    }

    load();
  }, [kernel.baseUrl, streamId]);

  // Kernel event updates
  useEffect(() => {
    const unsubscribe = eventBus.subscribe((event: KernelEvent) => {

      if (event.type === "issue.created") {
        if (event.payload.streamId !== streamId) return;

        setIssues(prev => {
          if (prev.some(i => i.id === event.payload.id)) return prev;
          return [...prev, event.payload];
        });
      }

      if (event.type === "issue.updated") {
        setIssues(prev =>
          prev.map(issue =>
            issue.id === event.payload.id ? event.payload : issue
          )
        );
      }

      if (event.type === "issue.deleted") {
        setIssues(prev =>
          prev.filter(issue => issue.id !== event.payload.id)
        );
      }

    });

    return unsubscribe;
  }, [streamId]);


  return issues;
}
