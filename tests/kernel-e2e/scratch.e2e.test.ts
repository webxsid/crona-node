import { describe, it, beforeAll, afterAll, expect } from "vitest";
import {
  startTestKernel,
  stopTestKernel,
  type IKernelTestHandle
} from "./helpers/kernel";
import { api } from "./helpers/http";

type ScratchPadRegisterResponse = {
  ok: boolean;
  filePath: string;
};

type ScratchpadMeta = {
  id: string;
  path: string;
  name: string;
  pinned?: boolean;
  lastOpenedAt?: string;
};

type ScratchpdCreateMeta = Omit<ScratchpadMeta, "id" | "lastOpenedAt" | "pinned">;

describe("@scratch @e2e", () => {
  let kernel: IKernelTestHandle;

  beforeAll(async () => {
    kernel = await startTestKernel();
  });

  afterAll(async () => {
    await stopTestKernel(kernel);
  });

  it("exposes scratch directory via kernel info", async () => {
    const info = await api<{ scratchDir: string }>(
      kernel.baseUrl,
      "/kernel/info",
      {
        headers: {
          Authorization: `Bearer ${kernel.token}`,
        },
      }
    );

    expect(info).toHaveProperty("scratchDir");
    expect(info.scratchDir).toContain(".crona");
  });

  it("registers a scratchpad", async () => {
    const scratch: ScratchpdCreateMeta = {
      path: "notes/today.md",
      name: "Today Notes",
    };

    const res = await api<ScratchPadRegisterResponse>(
      kernel.baseUrl,
      "/scratchpads/register",
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${kernel.token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(scratch),
      }
    );

    expect(res.ok).toBe(true);
    expect(res.filePath).toBe("notes/today.md");
  });

  it("lists scratchpads", async () => {
    const list = await api<ScratchpadMeta[]>(
      kernel.baseUrl,
      "/scratchpads",
      {
        headers: {
          Authorization: `Bearer ${kernel.token}`,
        },
      }
    );

    expect(Array.isArray(list)).toBe(true);
    expect(list.length).toBeGreaterThan(0);

    const entry = list.find(s => s.path === "notes/today.md");
    expect(entry).toBeDefined();
    expect(entry?.name).toBe("Today Notes");
  });

  it("pins and unpins a scratchpad", async () => {
    let list = await api<ScratchpadMeta[]>(
      kernel.baseUrl,
      "/scratchpads",
      {
        headers: {
          Authorization: `Bearer ${kernel.token}`,
        },
      }
    );

    console.log("Scratchpads before pinning:", list);

    const entry = list.find(s => s.path === "notes/today.md");
    expect(entry).toBeDefined();
    expect(entry?.pinned).toBeFalsy();

    const pinRes = await api<{ ok: boolean }>(
      kernel.baseUrl,
      `/scratchpads/pin/${entry?.id}`,
      {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${kernel.token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ pinned: true }),
      }
    );

    expect(pinRes.ok).toBe(true);

    list = await api<ScratchpadMeta[]>(
      kernel.baseUrl,
      "/scratchpads",
      {
        headers: {
          Authorization: `Bearer ${kernel.token}`,
        },
      }
    );

    expect(list.find(s => s.path === "notes/today.md")?.pinned).toBe(true);

    // unpin
    await api(
      kernel.baseUrl,
      `/scratchpads/pin/${entry?.id}`,
      {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${kernel.token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ pinned: false }),
      }
    );

    list = await api<ScratchpadMeta[]>(
      kernel.baseUrl,
      "/scratchpads",
      {
        headers: {
          Authorization: `Bearer ${kernel.token}`,
        },
      }
    );

    expect(list.find(s => s.path === "notes/today.md")?.pinned).toBe(false);
  });

  // read the registered scratchpad content (should have an # header with the name)
  it("reads scratchpad content", async () => {
    const list = await api<ScratchpadMeta[]>(
      kernel.baseUrl,
      "/scratchpads",
      {
        headers: {
          Authorization: `Bearer ${kernel.token}`,
        },
      }
    );

    const entry = list.find(s => s.path === "notes/today.md");
    expect(entry).toBeDefined();

    const res = await api<{ ok: boolean; content: string }>(
      kernel.baseUrl,
      `/scratchpads/read/${entry?.id}`,
      {
        headers: {
          Authorization: `Bearer ${kernel.token}`,
        },
      }
    );

    expect(res.ok).toBe(true);
    expect(res.content).toContain("# Today Notes");
  });

  /* -------------------------------------------------------------------------- */
  /*                            VARIABLE PATH TESTS                             */
  /* -------------------------------------------------------------------------- */


  it("expands date variable in scratchpad path", async () => {
    const res = await api<ScratchPadRegisterResponse>(
      kernel.baseUrl,
      "/scratchpads/register",
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${kernel.token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          path: "daily/[[date]].md",
          name: "Daily Notes",
        }),
      }
    );

    expect(res.ok).toBe(true);


    const today = new Date().toISOString().split("T")[0];
    expect(res.filePath).toBe(`daily/${today}.md`);

  });

  it("supports multiple variables in path", async () => {
    const res = await api<ScratchPadRegisterResponse>(
      kernel.baseUrl,
      "/scratchpads/register",
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${kernel.token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          path: "sessions/[[date]]/[[time]]-notes.md",
          name: "Session Notes",
        }),
      }
    );

    expect(res.ok).toBe(true);
    const now = new Date();
    const expectedPath = `sessions/${now.toISOString().split("T")[0]}/${now.toTimeString()?.split(" ")?.[0]?.replace(/:/g, "-")}-notes.md`;
    expect(res.filePath).toBe(expectedPath);
  });

  it("rejects unsupported variables in path", async () => {
    const res = await fetch(
      `${kernel.baseUrl}/scratchpads/register`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${kernel.token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          path: "notes/[[invalid]].md",
          name: "Invalid Notes",
        }),
      }
    );

    const body = await res.json();

    console.log(body);

    expect(res.status).toBe(500);
    expect(body.message).toContain("Invalid variable in path");
  });

});

