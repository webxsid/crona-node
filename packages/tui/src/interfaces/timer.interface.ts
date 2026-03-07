import type { SessionSegmentType } from "./session_segment.interface.js";

export type TimerStatePayload =
  | {
    state: "idle";
  }
  | {
    state: "running";
    sessionId: string;
    issueId: string;
    segmentType: SessionSegmentType;
    elapsedSeconds: number;
  }
  | {
    state: "paused";
    issueId: string;
    elapsedSeconds: number;
  };
