import type { KernelInfo } from "../interfaces/kernel.interface.js";
import type { Stream } from "../interfaces/stream.interface.js";

function headers(kernel: KernelInfo) {
  return {
    Authorization: `Bearer ${kernel.token}`,
    "Content-Type": "application/json",
  };
}

/**
 * GET /streams?repoId=...
 */
export async function getStreams(
  kernel: KernelInfo,
  repoId?: string
): Promise<Stream[]> {

  if (!repoId) {
    throw new Error("repoId is required to fetch streams");
  }

  const res = await fetch(
    `${kernel.baseUrl}/streams?repoId=${repoId}`,
    {
      headers: headers(kernel),
    }
  );

  if (!res.ok) {
    throw new Error(`Failed to fetch streams (${res.status})`);
  }

  return res.json();
}

/**
 * POST /stream
 */
export async function createStreamApi(
  kernel: KernelInfo,
  repoId: string,
  name: string
): Promise<Stream> {

  const res = await fetch(`${kernel.baseUrl}/stream`, {
    method: "POST",
    headers: headers(kernel),
    body: JSON.stringify({
      repoId,
      name,
    }),
  });

  if (!res.ok) {
    throw new Error("Failed to create stream");
  }

  return res.json();
}

/**
 * PUT /stream/:id
 */
export async function updateStreamApi(
  kernel: KernelInfo,
  id: string,
  name: string
): Promise<Stream> {

  const res = await fetch(`${kernel.baseUrl}/stream/${id}`, {
    method: "PUT",
    headers: headers(kernel),
    body: JSON.stringify({ name }),
  });

  if (!res.ok) {
    throw new Error("Failed to update stream");
  }

  return res.json();
}

/**
 * DELETE /stream/:id
 */
export async function deleteStreamApi(
  kernel: KernelInfo,
  id: string
): Promise<void> {

  const res = await fetch(`${kernel.baseUrl}/stream/${id}`, {
    method: "DELETE",
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to delete stream");
  }
}
