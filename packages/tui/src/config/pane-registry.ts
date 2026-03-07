export type Pane =
  | "repos"
  | "streams"
  | "issues"
  | "scratchpads"
  | "command";

export type PaneActionId =
  | "repo.create"
  | "repo.delete"
  | "stream.create"
  | "stream.delete"
  | "issue.create"
  | "issue.delete"
  | "issue.todo.toggle"
  | "scratchpad.create"
  | "scratchpad.open";

export interface PaneActionDefinition {
  key: string;
  id: PaneActionId;
  label: string;
}

export interface PaneDefinition {
  title: string;
  actions: PaneActionDefinition[];
}

export const paneRegistry: Record<Exclude<Pane, "command">, PaneDefinition> = {
  repos: {
    title: "Repos",
    actions: [
      { key: "a", id: "repo.create", label: "add repo" },
      { key: "d", id: "repo.delete", label: "delete repo" },
    ],
  },

  streams: {
    title: "Streams",
    actions: [
      { key: "a", id: "stream.create", label: "add stream" },
      { key: "d", id: "stream.delete", label: "delete stream" },
    ],
  },

  issues: {
    title: "Issues",
    actions: [
      { key: "a", id: "issue.create", label: "add issue" },
      { key: "t", id: "issue.todo.toggle", label: "toggle today" },
      { key: "d", id: "issue.delete", label: "delete issue" },
    ],
  },

  scratchpads: {
    title: "Scratchpads",
    actions: [
      { key: "a", id: "scratchpad.create", label: "add scratchpad" },
      { key: "o", id: "scratchpad.open", label: "open scratchpad" },
    ],
  },
};


export function getPaneActions(pane: Pane): PaneActionDefinition[] {
  if (pane === "command") return [];
  return paneRegistry[pane].actions;
}

export function getPaneActionByKey(
  pane: Pane,
  key: string
): PaneActionDefinition | undefined {
  if (pane === "command") return undefined;
  return paneRegistry[pane].actions.find((action) => action.key === key);
}
