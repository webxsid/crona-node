import type { Session } from "node:sqlite";
import type { Issue } from "./issue.interface.js";
import type { Repo } from "./repo.interface.js";
import type { Stream } from "./stream.interface.js";
import type { TimerStatePayload } from "./timer.interface.js";
import type { ActiveContext } from "./context.interface.js";
import type { SessionSegmentType } from "./session_segment.interface.js";
import type { ScratchPadMeta } from "./scratchpad_meta.interface.js";
import type { Stash } from "./stash.interface.js";

export type KernelEvent =
  | { type: "repo.created"; payload: Repo }
  | { type: "repo.updated"; payload: Repo }
  | { type: "repo.deleted"; payload: { id: string } }
  | { type: "stream.created"; payload: Stream }
  | { type: "stream.updated"; payload: Stream }
  | { type: "stream.deleted"; payload: { id: string } }
  | { type: "issue.created"; payload: Issue }
  | { type: "issue.updated"; payload: Issue }
  | { type: "issue.deleted"; payload: { id: string } }
  | { type: "session.started"; payload: Session }
  | { type: "session.stopped"; payload: Session }
  | { type: "timer.state"; payload: TimerStatePayload }
  | { type: "context.changed"; payload: Pick<ActiveContext, "deviceId" | "repoId" | "streamId" | "issueId"> }
  | { type: "stash.created" | "stash.applied"; payload: Pick<Stash, "id" | "deviceId" | "repoId" | "streamId" | "issueId"> }
  | { type: "stash.dropped"; payload: { id: string; deviceId: string } }
  | { type: "timer.boundary"; payload: { from: SessionSegmentType; to: SessionSegmentType } }
  | { type: "timer.tick"; payload: { remainingSeconds: number } }
  | { type: "scratchpad.created"; payload: ScratchPadMeta }
  | { type: "scratchpad.updated"; payload: ScratchPadMeta }
  | { type: "scratchpad.deleted"; payload: { id: string } }
  ;

