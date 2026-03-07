import React, { useEffect, useState } from "react";
import { Text } from "ink";
import { App } from "./app.js";
import { ensureKernelRunning } from "./kernel/launcher.js";
import { watchKernel } from "./kernel/watchdog.js";
import { clearTerminal } from "./utils/terminal.js";
import { KernelInfo } from "./interfaces/kernel.interface.js";



interface BootstrapProps {
  unmount: () => void;
}

let stopWatch: (() => void) | null = null;

export function Bootstrap({ unmount }: BootstrapProps) {
  const [kernel, setKernel] = useState<KernelInfo | null>(null);

  useEffect(() => {
    ensureKernelRunning()
      .then(setKernel)
      .catch((err) => {
        console.error("Failed to start kernel:", err);
        process.exit(1);
      });
  }, []);

  useEffect(() => {
    if (!kernel) return;

    stopWatch = watchKernel(kernel, (reason) => {
      unmount();
      clearTerminal();

      console.error(
        "\ncrona: kernel disconnected\n" +
        `reason: ${reason}\n`
      );

      process.exit(1);
    });

    return () => {
      stopWatch?.();
      stopWatch = null;
    };
  }, [kernel, unmount]);

  if (!kernel) {
    return <Text>Starting Crona kernel…</Text>;
  }

  return <App kernel={kernel} />;
}
