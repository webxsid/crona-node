import fs from "node:fs/promises";
import path from "node:path";
import { CRONA_SCRATCH_DIR } from "./kernel-info";

/**
 * Resolve scratchpad path safely inside the scratch directory
 */
function resolveScratchPath(userPath: string): string {
  if (!userPath) {
    throw new Error("Scratchpad path is required");
  }

  let normalized = userPath;

  if (!normalized.endsWith(".md")) {
    normalized += ".md";
  }

  const fullPath = path.resolve(CRONA_SCRATCH_DIR, normalized);

  // prevent directory escape
  const root = path.resolve(CRONA_SCRATCH_DIR);
  if (!fullPath.startsWith(root)) {
    throw new Error("Invalid scratchpad path");
  }

  return fullPath;
}

/**
 * Create scratchpad file if it does not exist
 */
export async function createScratchPadFile(
  userPath: string,
  title?: string
): Promise<string> {
  const fullPath = resolveScratchPath(userPath);

  const dir = path.dirname(fullPath);

  await fs.mkdir(dir, { recursive: true });

  try {
    const handle = await fs.open(fullPath, "wx");

    const safeTitle = title?.trim();
    if (safeTitle) {
      await handle.writeFile(`# ${safeTitle}\n\n`, "utf8");
    }

    await handle.close();
  } catch (err) {
    if ((err as NodeJS.ErrnoException).code !== "EEXIST") {
      throw err;
    }
  }

  return fullPath;
}

export async function readScratchPadFile(userPath: string): Promise<string> {
  const fullPath = resolveScratchPath(userPath);

  try {
    return await fs.readFile(fullPath, "utf8");
  } catch (err) {
    if ((err as NodeJS.ErrnoException).code === "ENOENT") {
      throw new Error("Scratchpad file not found");
    }
    throw err;
  }
}

/**
 * Delete scratchpad file (optional helper)
 * Metadata deletion happens separately.
 */
export async function deleteScratchPadFile(userPath: string): Promise<void> {
  const fullPath = resolveScratchPath(userPath);

  try {
    await fs.unlink(fullPath);
  } catch (err) {
    if ((err as NodeJS.ErrnoException).code !== "ENOENT") {
      throw err;
    }
  }
}
