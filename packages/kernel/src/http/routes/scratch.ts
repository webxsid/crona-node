import type { FastifyInstance } from "fastify";
import type { ScratchPadMeta } from "@crona/core";
import { getScratchpad, listScratchpads, pinScratchpad, registerScratchpad, removeScratchpad, type ICommandContext } from "@crona/core";
import { createScratchPadFile, readScratchPadFile } from "../../scratch-file";

export class ScratchRoutes {
  constructor(
    private readonly app: FastifyInstance,
    private readonly ctx: ICommandContext
  ) { }

  register() {
    this.registerQueries();
    this.registerCommands();
  }

  private registerQueries() {
    // GET /scratchpads
    this.app.get("/scratchpads", async (req) => {
      const pinnedOnly = (req.query as { pinnedOnly: string }).pinnedOnly === "1";
      return listScratchpads(this.ctx, {
        pinnedOnly
      });
    });
  }

  private registerCommands() {
    // POST /scratchpads/register
    this.app.post("/scratchpads/register", async (req) => {
      const body = req.body as ScratchPadMeta;

      // core command (DB metadata)
      const filePath = await registerScratchpad(this.ctx, body);

      try {
        // kernel filesystem responsibility
        await createScratchPadFile(filePath, body.name);
      } catch (err) {
        // rollback metadata if file creation fails
        await removeScratchpad(this.ctx, filePath);
        throw err;
      }

      return { ok: true, filePath };
    });

    this.app.get("/scratchpads/meta/:id", async (req) => {
      const { id } = req.params as { id: string };
      const meta = await getScratchpad(this.ctx, id);
      if (!meta) {
        return { ok: false, error: "Scratchpad not found" };
      }
      return { ok: true, meta };
    });


    this.app.get("/scratchpads/read/:id", async (req) => {
      const { id } = req.params as { id: string };
      const meta = await getScratchpad(this.ctx, id);
      if (!meta) {
        return { ok: false, error: "Scratchpad not found" };
      }

      const fileData = await readScratchPadFile(meta.path)

      return { ok: true, meta, content: fileData };
    });

    // PUT /scratchpads/pin/:id.  body: { pinned: boolean }
    this.app.put("/scratchpads/pin/:id", async (req) => {
      const { id } = req.params as { id: string };
      const { pinned } = req.body as { pinned: boolean };
      await pinScratchpad(this.ctx, id, pinned);
      return { ok: true };
    });

    // DELETE /scratchpads/:id.
    this.app.delete("/scratchpads/:id", async (req) => {
      const { id } = req.params as { id: string };
      await removeScratchpad(this.ctx, id);
      return { ok: true };
    });
  }
}
