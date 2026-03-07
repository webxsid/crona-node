import type { KernelInfo } from "../interfaces/kernel.interface.js";
import type { ScratchPadMeta } from "../interfaces/scratchpad_meta.interface.js";

function headers(kernel: KernelInfo) {
  return {
    Authorization: `Bearer ${kernel.token}`,
    "Content-Type": "application/json",
  };
}

/**
 * GET /scratchpads
 */
export async function getScratchpads(
  kernel: KernelInfo,
  pinnedOnly = false
): Promise<ScratchPadMeta[]> {

  const url = pinnedOnly
    ? `${kernel.baseUrl}/scratchpads?pinnedOnly=1`
    : `${kernel.baseUrl}/scratchpads`;

  const res = await fetch(url, {
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to fetch scratchpads");
  }

  return res.json();
}

/**
 * POST /scratchpads/register
 */
export async function createScratchpad(
  kernel: KernelInfo,
  meta: ScratchPadMeta
) {

  const res = await fetch(`${kernel.baseUrl}/scratchpads/register`, {
    method: "POST",
    headers: headers(kernel),
    body: JSON.stringify(meta),
  });

  if (!res.ok) {
    throw new Error("Failed to create scratchpad");
  }

  return res.json();
}

/**
 * GET /scratchpads/read/:id
 */
export async function readScratchpad(
  kernel: KernelInfo,
  id: string
) {
  const res = await fetch(`${kernel.baseUrl}/scratchpads/read/${id}`, {
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to read scratchpad");
  }

  return res.json();
}

/**
 * PUT /scratchpads/pin/:id
 */
export async function pinScratchpadApi(
  kernel: KernelInfo,
  id: string,
  pinned: boolean
) {
  const res = await fetch(`${kernel.baseUrl}/scratchpads/pin/${id}`, {
    method: "PUT",
    headers: headers(kernel),
    body: JSON.stringify({ pinned }),
  });

  if (!res.ok) {
    throw new Error("Failed to pin scratchpad");
  }
}

/**
 * DELETE /scratchpads/:id
 */
export async function deleteScratchpad(
  kernel: KernelInfo,
  id: string
) {
  const res = await fetch(`${kernel.baseUrl}/scratchpads/${id}`, {
    method: "DELETE",
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to delete scratchpad");
  }
}
