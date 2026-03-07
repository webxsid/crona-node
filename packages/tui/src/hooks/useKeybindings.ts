import { useInput } from "ink";
import type { KernelInfo } from "../interfaces/kernel.interface.js";
import type { Pane } from "./useTuiState.js";
import {
  startTimer,
  pauseTimer,
  resumeTimer,
  endTimer
} from "../services/timer-api.js";

export interface KeybindingActions {
  pane: Pane;

  moveUp: () => void;
  moveDown: (max: number) => void;

  focusNextPane: () => void;
  focusPrevPane: () => void;
  focusPane: (pane: Pane) => void;

  setCommandBuffer: (v: string) => void;
  clearCommand: () => void;
}

export function useKeybindings(
  kernel: KernelInfo,
  actions: KeybindingActions,
  getPaneMax: () => number
) {

  useInput(async (input, key) => {

    // ---------- pane navigation ----------

    if (key.tab) {
      actions.focusNextPane();
      return;
    }

    if (key.shift && key.tab) {
      actions.focusPrevPane();
      return;
    }

    // direct pane jumps
    if (input === "1") actions.focusPane("repos");
    if (input === "2") actions.focusPane("streams");
    if (input === "3") actions.focusPane("issues");
    if (input === "4") actions.focusPane("scratchpads");

    // ---------- movement ----------

    if (key.upArrow || input === "k") {
      actions.moveUp();
      return;
    }

    if (key.downArrow || input === "j") {
      actions.moveDown(getPaneMax());
      return;
    }

    // ---------- command mode ----------

    if (input === ":") {
      actions.focusPane("command");
      actions.clearCommand();
      return;
    }

    if (actions.pane === "command") {

      if (key.escape) {
        actions.focusPane("issues");
        actions.clearCommand();
        return;
      }

      if (key.return) {
        // handled later by useCommand
        return;
      }

      if (key.backspace || key.delete) {
        actions.setCommandBuffer("");
        return;
      }

      actions.setCommandBuffer(input);
      return;
    }

    // ---------- timer controls ----------

    if (input === "s") {
      await startTimer(kernel);
      return;
    }

    if (input === "p") {
      await pauseTimer(kernel);
      return;
    }

    if (input === "r") {
      await resumeTimer(kernel);
      return;
    }

    if (input === "e") {
      await endTimer(kernel);
      return;
    }

  });
}
