import type { Stream, StreamVisibility } from "../../domain";
import { SqliteDb } from "../../storage";
import type { IStreamRepository } from "./interface";

export class SqliteStreamRepository implements IStreamRepository {
  async create(
    stream: Stream,
    meta: { userId: string; now: string }
  ): Promise<Stream> {
    await SqliteDb.getDB()
      .insertInto("streams")
      .values({
        id: stream.id,
        repo_id: stream.repoId,
        name: stream.name,
        visibility: stream.visibility,
        user_id: meta.userId,
        created_at: meta.now,
        updated_at: meta.now,
        deleted_at: null,
      })
      .execute();

    return stream;
  }

  async getById(
    streamId: string,
    userId: string
  ): Promise<Stream | null> {
    const row = await SqliteDb.getDB()
      .selectFrom("streams")
      .select(["id", "repo_id", "name", "visibility"])
      .where("id", "=", streamId)
      .where("user_id", "=", userId)
      .where("deleted_at", "is", null)
      .executeTakeFirst();

    if (!row) return null;

    return {
      id: row.id,
      repoId: row.repo_id,
      name: row.name,
      visibility: row.visibility as StreamVisibility,
    };
  }

  async listByRepo(
    repoId: string,
    userId: string
  ): Promise<Stream[]> {
    const rows = await SqliteDb.getDB()
      .selectFrom("streams")
      .select(["id", "repo_id", "name", "visibility"])
      .where("repo_id", "=", repoId)
      .where("user_id", "=", userId)
      .where("deleted_at", "is", null)
      .orderBy("created_at", "asc")
      .execute();

    return rows.map((r) => ({
      id: r.id,
      repoId: r.repo_id,
      name: r.name,
      visibility: r.visibility as StreamVisibility,
    }));
  }

  async listDeletedByRepo(
    repoId: string,
    userId: string
  ): Promise<Stream[]> {
    // join repo to get streams that are deleted along with the repo which is not deleted, to show in the UI as deleted streams
    const rows = await SqliteDb.getDB()
      .selectFrom("streams")
      .innerJoin("repos", "repos.id", "streams.repo_id")
      .select(["streams.id", "streams.repo_id", "streams.name", "streams.visibility"])
      .where("streams.repo_id", "=", repoId)
      .where("streams.user_id", "=", userId)
      .where("streams.deleted_at", "is not", null)
      .where("repos.deleted_at", "is", null)
      .orderBy("streams.created_at", "asc")
      .execute();

    return rows.map((r) => ({
      id: r.id,
      repoId: r.repo_id,
      name: r.name,
      visibility: r.visibility as StreamVisibility,
    }));
  }

  async update(
    streamId: string,
    updates: {
      name?: string;
      visibility?: StreamVisibility;
    },
    meta: { userId: string; now: string }
  ): Promise<Stream> {
    const result = await SqliteDb.getDB()
      .updateTable("streams")
      .set({
        name: updates.name,
        visibility: updates.visibility,
        updated_at: meta.now,
      })
      .where("id", "=", streamId)
      .where("user_id", "=", meta.userId)
      .where("deleted_at", "is", null)
      .executeTakeFirst();

    if (result.numUpdatedRows === BigInt(0)) {
      throw new Error("Stream not found or already deleted");
    }

    const updated = await this.getById(streamId, meta.userId);
    if (!updated) {
      throw new Error("Stream disappeared after update");
    }

    return updated;
  }

  async softDelete(
    streamId: string,
    meta: { userId: string; now: string }
  ): Promise<void> {
    const result = await SqliteDb.getDB()
      .updateTable("streams")
      .set({
        deleted_at: meta.now,
        updated_at: meta.now,
      })
      .where("id", "=", streamId)
      .where("user_id", "=", meta.userId)
      .where("deleted_at", "is", null)
      .executeTakeFirst();

    if (result.numUpdatedRows === BigInt(0)) {
      throw new Error("Stream not found or already deleted");
    }
  }

  async restore(
    streamId: string,
    meta: { userId: string; now: string }
  ): Promise<void> {
    const result = await SqliteDb.getDB()
      .updateTable("streams")
      .set({
        deleted_at: null,
        updated_at: meta.now,
      })
      .where("id", "=", streamId)
      .where("user_id", "=", meta.userId)
      .where("deleted_at", "is not", null)
      .executeTakeFirst();

    if (result.numUpdatedRows === BigInt(0)) {
      throw new Error("Stream not found or not deleted");
    }
  }

  async cascadeSoftDeleteByRepoId(repoId: string, meta: { userId: string; now: string; }): Promise<void> {
    await SqliteDb.getDB()
      .updateTable("streams")
      .set({
        deleted_at: meta.now,
        updated_at: meta.now,
      })
      .where("repo_id", "=", repoId)
      .where("user_id", "=", meta.userId)
      .where("deleted_at", "is", null)
      .execute();
  }

  async restoreByRepoId(repoId: string, meta: { userId: string; now: string; }): Promise<void> {
    await SqliteDb.getDB()
      .updateTable("streams")
      .set({
        deleted_at: null,
        updated_at: meta.now,
      })
      .where("repo_id", "=", repoId)
      .where("user_id", "=", meta.userId)
      .where("deleted_at", "is not", null)
      .execute();
  }
}
