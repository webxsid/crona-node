import { randomUUID } from "crypto";
import type { Stream, StreamVisibility } from "../domain/stream";
import type { ICommandContext } from "./context";
import { cacadeSoftDeleteIssuesByStreamId } from "./issue.commands";

/**
 * Create a new stream under a repo
 */
export async function createStream(
  ctx: ICommandContext,
  input: {
    repoId: string;
    name: string;
    visibility?: StreamVisibility;
  }
): Promise<Stream> {
  if (!input.name.trim()) {
    throw new Error("Stream name cannot be empty");
  }

  const stream: Stream = {
    id: randomUUID(),
    repoId: input.repoId,
    name: input.name.trim(),
    visibility: input.visibility ?? "personal",
  };

  const now = ctx.now();

  await ctx.streams.create(stream, {
    userId: ctx.userId,
    now,
  });

  await ctx.ops.append({
    id: randomUUID(),
    entity: "stream",
    entityId: stream.id,
    action: "create",
    payload: stream,
    timestamp: now,
    userId: ctx.userId,
    deviceId: ctx.deviceId,
  });

  return stream;
}
/**
 * Rename / change visibility of a stream
 */
export async function updateStream(
  ctx: ICommandContext,
  streamId: string,
  updates: {
    name?: string;
    visibility?: StreamVisibility;
  }
): Promise<Stream> {
  if (updates.name !== undefined && !updates.name.trim()) {
    throw new Error("Stream name cannot be empty");
  }

  const now = ctx.now();

  const updated = await ctx.streams.update(
    streamId,
    {
      name: updates.name?.trim(),
      visibility: updates.visibility,
    },
    {
      userId: ctx.userId,
      now,
    }
  );

  await ctx.ops.append({
    id: randomUUID(),
    entity: "stream",
    entityId: streamId,
    action: "update",
    payload: updates,
    timestamp: now,
    userId: ctx.userId,
    deviceId: ctx.deviceId,
  });

  return updated;
}
/**
 * Delete a stream (soft delete)
 */
export async function deleteStream(
  ctx: ICommandContext,
  streamId: string
): Promise<void> {
  const now = ctx.now();

  await ctx.streams.softDelete(streamId, {
    userId: ctx.userId,
    now,
  });

  await ctx.ops.append({
    id: randomUUID(),
    entity: "stream",
    entityId: streamId,
    action: "delete",
    payload: null,
    timestamp: now,
    userId: ctx.userId,
    deviceId: ctx.deviceId,
  });

  await cacadeSoftDeleteIssuesByStreamId(ctx, streamId);
}

export async function restoreStream(
  ctx: ICommandContext,
  streamId: string
): Promise<void> {
  const now = ctx.now();

  await ctx.streams.restore(streamId, {
    userId: ctx.userId,
    now,
  });

  await ctx.ops.append({
    id: randomUUID(),
    entity: "stream",
    entityId: streamId,
    action: "restore",
    payload: null,
    timestamp: now,
    userId: ctx.userId,
    deviceId: ctx.deviceId,
  });

}

/**
 * Cascade delete all streams under a repo (soft delete)
 */

export async function cascadeDeleteStreams(
  ctx: ICommandContext,
  repoId: string,
): Promise<void> {
  const now = ctx.now();

  await ctx.streams.cascadeSoftDeleteByRepoId(repoId, {
    userId: ctx.userId,
    now,
  });
}

export async function cascadeRestoreStreams(
  ctx: ICommandContext,
  repoId: string,
): Promise<void> {
  const now = ctx.now();

  await ctx.streams.restoreByRepoId(repoId, {
    userId: ctx.userId,
    now,
  });
}

/**
 * List streams for a repo
 * Read-only → no ops
 */
export async function listStreamsByRepo(
  ctx: ICommandContext,
  repoId: string
): Promise<Stream[]> {
  return ctx.streams.listByRepo(repoId, ctx.userId);
}
