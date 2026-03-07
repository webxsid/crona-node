import { useEffect, useState } from "react";
import type { KernelEvent } from "../interfaces/events.interface.js";
import { eventBus } from "../state/event-bus.js";

export type Pane =
  | "repos"
  | "streams"
  | "issues"
  | "scratchpads"
  | "command";

export interface TuiState {
  pane: Pane;

  cursor: {
    repos: number;
    streams: number;
    issues: number;
    scratchpads: number;
  };

  commandBuffer: string;
}

const paneOrder: Pane[] = [
  "repos",
  "streams",
  "issues",
  "scratchpads",
  "command"
];

export function useTuiState() {
  const [state, setState] = useState<TuiState>({
    pane: "issues",

    cursor: {
      repos: 0,
      streams: 0,
      issues: 0,
      scratchpads: 0
    },

    commandBuffer: ""
  });

  // ---------- Pane focus ----------

  function focusPane(pane: Pane) {
    setState((s) => ({ ...s, pane }));
  }

  function focusNextPane() {
    setState((s) => {
      const index = paneOrder.indexOf(s.pane);
      const next = paneOrder[(index + 1) % paneOrder.length];

      return { ...s, pane: next };
    });
  }

  function focusPrevPane() {
    setState((s) => {
      const index = paneOrder.indexOf(s.pane);
      const prev =
        paneOrder[(index - 1 + paneOrder.length) % paneOrder.length];

      return { ...s, pane: prev };
    });
  }

  // ---------- Cursor movement ----------

  function moveUp() {
    setState((s) => {
      const pane = s.pane;
      if (pane === "command") return s;

      const key = pane as keyof TuiState["cursor"];

      return {
        ...s,
        cursor: {
          ...s.cursor,
          [key]: Math.max(0, s.cursor[key] - 1)
        }
      };
    });
  }

  function moveDown(max: number) {
    setState((s) => {
      const pane = s.pane;
      if (pane === "command") return s;

      const key = pane as keyof TuiState["cursor"];

      return {
        ...s,
        cursor: {
          ...s.cursor,
          [key]: Math.min(max - 1, s.cursor[key] + 1)
        }
      };
    });
  }

  // ---------- Command buffer ----------

  function setCommandBuffer(value: string) {
    setState((s) => ({
      ...s,
      commandBuffer: value
    }));
  }

  function clearCommand() {
    setState((s) => ({
      ...s,
      commandBuffer: ""
    }));
  }

  // ---------- Context change handling ----------

  useEffect(() => {
    const unsubscribe = eventBus.subscribe((event: KernelEvent) => {

      if (event.type !== "context.changed") return;

      const { repoId, streamId, issueId } = event.payload;

      setState((s) => {
        const next = { ...s };

        // repo changed → reset streams + issues
        if (!repoId) {
          next.cursor.streams = 0;
          next.cursor.issues = 0;
        }

        // stream changed → reset issues
        if (!streamId) {
          next.cursor.issues = 0;
        }

        // issue cleared → reset issue cursor
        if (!issueId) {
          next.cursor.issues = 0;
        }

        return next;
      });
    });

    return unsubscribe;
  }, []);

  return {
    state,

    focusPane,
    focusNextPane,
    focusPrevPane,

    moveUp,
    moveDown,

    setCommandBuffer,
    clearCommand
  };
}
