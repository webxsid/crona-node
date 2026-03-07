import { useEffect, useState } from "react";
import type { KernelInfo } from "../interfaces/kernel.interface.js";
import type { Repo } from "../interfaces/repo.interface.js";
import { getRepos } from "../services/repo-api.js";
import type { KernelEvent } from "../interfaces/events.interface.js";
import { eventBus } from "../state/event-bus.js";

export function useRepos(kernel: KernelInfo) {
  const [repos, setRepos] = useState<Repo[]>([]);

  useEffect(() => {
    async function load() {
      try {
        const data = await getRepos(kernel);
        setRepos(data);
      } catch (err) {
        console.error("Failed to load repos", err);
      }
    }

    load();
  }, [kernel.baseUrl]);

  // Live kernel updates
  useEffect(() => {
    const unsubscribe = eventBus.subscribe((event: KernelEvent) => {
      if (event.type === "repo.created") {
        setRepos(prev => {
          if (prev.some(r => r.id === event.payload.id)) return prev;
          return [...prev, event.payload];
        });
      }

      if (event.type === "repo.updated") {
        setRepos(prev =>
          prev.map(repo =>
            repo.id === event.payload.id ? event.payload : repo
          )
        );
      }

      if (event.type === "repo.deleted") {
        setRepos(prev =>
          prev.filter(repo => repo.id !== event.payload.id)
        );
      }
    });

    return unsubscribe;
  }, []);

  return repos;
}
