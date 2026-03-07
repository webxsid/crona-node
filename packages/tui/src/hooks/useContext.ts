import { useEffect, useState } from "react";
import type { KernelInfo } from "../interfaces/kernel.interface.js";
import type { ActiveContext } from "../interfaces/context.interface.js";
import { getContext } from "../services/context-api.js";

export function useContext(kernel: KernelInfo) {
  const [context, setContext] = useState<ActiveContext | null>(null);

  useEffect(() => {
    async function load() {
      try {
        const ctx = await getContext(kernel);
        setContext(ctx);
      } catch { }
    }

    load();
  }, [kernel]);

  return context;
}
