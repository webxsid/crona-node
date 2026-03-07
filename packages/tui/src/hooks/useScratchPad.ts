import { useEffect, useState } from "react";
import type { KernelInfo } from "../interfaces/kernel.interface.js";
import { readScratchpad } from "../services/scratchpad-api.js";

export function useScratchpad(
  kernel: KernelInfo,
  id?: string
) {

  const [content, setContent] = useState<string>("");

  useEffect(() => {

    if (!id) {
      setContent("");
      return;
    }

    async function load() {
      try {
        const res = await readScratchpad(kernel, id!);
        setContent(res.content);
      } catch (err) {
        console.error("Failed to load scratchpad", err);
      }
    }

    load();

  }, [kernel.baseUrl, id]);

  return content;
}
