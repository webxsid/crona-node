import { useEffect, useState } from "react";
import type { KernelInfo } from "../interfaces/kernel.interface.js";
import { getTimerState } from "../services/timer-api.js";
import type { TimerStatePayload } from "../interfaces/timer.interface.js";
import type { KernelEvent } from "../interfaces/events.interface.js";
import { eventBus } from "../state/event-bus.js";

export function useTimer(kernel: KernelInfo) {
  const [timer, setTimer] = useState<TimerStatePayload | null>(null);
  const [elapsedSeconds, setElapsedSeconds] = useState(0);

  // initial load
  useEffect(() => {
    let alive = true;

    async function load() {
      try {
        const state = await getTimerState(kernel);
        if (alive) {
          setTimer(state);
          setElapsedSeconds(0);
        }
      } catch (err) {
        console.error("Failed to fetch timer state", err);
      }
    }

    load();

    return () => {
      alive = false;
    };
  }, [kernel.baseUrl]);

  // local ticking
  useEffect(() => {
    if (!timer || timer.state === "idle") return;

    const interval = setInterval(() => {
      setElapsedSeconds((s) => s + 1);
    }, 1000);

    return () => clearInterval(interval);
  }, [timer?.state]);

  // kernel events
  useEffect(() => {
    const unsubscribe = eventBus.subscribe((event: KernelEvent) => {
      if (event.type === "timer.state") {
        setTimer(event.payload);
        setElapsedSeconds(0);
        return;
      }

      if (event.type === "timer.boundary") {
        setElapsedSeconds(0);
      }
    });

    return unsubscribe;
  }, []);

  return { timer, elapsedSeconds };
}
