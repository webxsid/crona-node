import type { ScratchPadMeta } from "../domain";
import type { ICommandContext } from "./context";
import * as path from "path";
import { v4 as uuidv4 } from "uuid";

export function validatePath(inputPath: string): void {
  if (!inputPath || inputPath.trim() === "") {
    throw new Error("Path cannot be empty");
  }

  // Normalize first
  const normalized = path.posix.normalize(inputPath);

  // Prevent absolute paths
  if (normalized.startsWith("/")) {
    throw new Error("Absolute paths are not allowed");
  }

  // Prevent traversal
  if (normalized.startsWith("../") || normalized.includes("/../")) {
    throw new Error("Path traversal is not allowed");
  }

  const invalidChars = /[<>:"\\|?*]/;

  const segments = normalized.split("/");

  for (const segment of segments) {
    if (!segment) continue;

    if (segment === "." || segment === "..") {
      throw new Error("Invalid path segment");
    }

    if (invalidChars.test(segment)) {
      throw new Error(
        "Path contains invalid characters: < > : \" \\ | ? *"
      );
    }
  }
}

function handleVariablesInPath(path: string): string {
  // the vriables will be in the format [[variableName]]
  // supported variables are:
  // - [[date]]: current date in YYYY-MM-DD format
  // - [[time]]: current time in HH-mm-ss format
  // - [[datetime]]: current date and time in YYYY-MM-DD_HH-mm-ss format
  // - [[timestamp]]: current timestamp in milliseconds
  // - [[random]]: random string of 8 characters

  const allowedVariables = new Set(["date", "time", "datetime", "timestamp", "random"]);

  const varibaleRegex = /\[\[(\w+)\]\]/g;

  const foundVariables = [...path.matchAll(varibaleRegex)];

  if (foundVariables.some((match) => {
    const variableName = match[1];
    return !allowedVariables.has(variableName!);
  })) {
    console.error("Invalid variable found in path:", Array.from(foundVariables));
    throw new Error(
      `Invalid variable in path. Allowed variables are: ${allowedVariables.size > 0 ? Array.from(allowedVariables).join(", ") : "none"}`
    );
  }

  let resultPath = path;

  for (const match of path.matchAll(varibaleRegex)) {
    const variableName = match[1];
    let replacement = "";
    const now = new Date();
    switch (variableName) {
      case "date":
        replacement = now.toISOString().split("T")[0]!;
        break;
      case "time":
        replacement = now.toTimeString().split(" ")[0]!.replace(/:/g, "-");
        break;
      case "datetime":
        replacement = now
          .toISOString()
          .replace(/[:.]/g, "-")
          .replace("T", "_")
          .split("Z")[0]!;
        break;
      case "timestamp":
        replacement = now.getTime().toString();
        break;
      case "random":
        replacement = Math.random().toString(36).substring(2, 10);
        break;
    }
    resultPath = resultPath.replace(match[0], replacement);
  }

  return resultPath;

}

export async function registerScratchpad(
  ctx: ICommandContext,
  meta: ScratchPadMeta
): Promise<string> {
  const incomingPath = meta.path;
  if (!incomingPath) {
    throw new Error("Path is required to register a scratchpad");
  }

  const normalizedPath = incomingPath.trim();
  validatePath(normalizedPath);

  const processedPath = handleVariablesInPath(normalizedPath);
  console.log(`Processed scratchpad path: ${processedPath}`);
  await ctx.scratchPads.upsert({
    id: uuidv4(),
    name: meta.name,
    path: processedPath,
    pinned: meta.pinned ?? false,
    lastOpenedAt: new Date(ctx.now()),
  }, {
    userId: ctx.userId,
    deviceId: ctx.deviceId,
  });

  return processedPath;
}

export async function getScratchpad(
  ctx: ICommandContext,
  id: string
): Promise<ScratchPadMeta | null> {
  return ctx.scratchPads.getById(id, {
    userId: ctx.userId,
    deviceId: ctx.deviceId,
  });
};

export async function listScratchpads(
  ctx: ICommandContext,
  options: {
    pinnedOnly?: boolean;
  }
): Promise<ScratchPadMeta[]> {
  return ctx.scratchPads.list(
    ctx.userId,
    ctx.deviceId,
    {
      pinnedOnly: options.pinnedOnly ?? false,
    }
  );
}

export async function pinScratchpad(
  ctx: ICommandContext,
  id: string,
  pinned: boolean
): Promise<void> {
  const existing = await ctx.scratchPads.getById(id, {
    userId: ctx.userId,
    deviceId: ctx.deviceId,
  });
  if (!existing) throw new Error("Scratchpad not found");

  await ctx.scratchPads.upsert({ ...existing, pinned }, {
    userId: ctx.userId,
    deviceId: ctx.deviceId,
  });
}

export async function removeScratchpad(
  ctx: ICommandContext,
  id: string
): Promise<void> {
  await ctx.scratchPads.removeBYId(id, {
    userId: ctx.userId,
    deviceId: ctx.deviceId,
  });
}
