import type { KernelInfo } from "../interfaces/kernel.interface.js";
import type { ActiveContext } from "../interfaces/context.interface.js";


function headers(kernel: KernelInfo) {
  return {
    Authorization: `Bearer ${kernel.token}`,
    "Content-Type": "application/json",
  };
}

/**
 * GET /context
 */
export async function getContext(kernel: KernelInfo): Promise<ActiveContext> {
  const res = await fetch(`${kernel.baseUrl}/context`, {
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error(`Failed to fetch context (${res.status})`);
  }

  return res.json();
}

/**
 * PUT /context/repo
 */
export async function switchRepo(
  kernel: KernelInfo,
  repoId: string
): Promise<ActiveContext> {
  const res = await fetch(
    `${kernel.baseUrl}/context/repo?repoId=${repoId}`,
    {
      method: "PUT",
      headers: headers(kernel),
    }
  );

  if (!res.ok) {
    throw new Error("Failed to switch repo");
  }

  return res.json();
}

/**
 * PUT /context/stream
 */
export async function switchStream(
  kernel: KernelInfo,
  streamId: string
): Promise<ActiveContext> {
  const res = await fetch(
    `${kernel.baseUrl}/context/stream?streamId=${streamId}`,
    {
      method: "PUT",
      headers: headers(kernel),
    }
  );

  if (!res.ok) {
    throw new Error("Failed to switch stream");
  }

  return res.json();
}

/**
 * PUT /context/issue
 */
export async function switchIssue(
  kernel: KernelInfo,
  issueId: string
): Promise<ActiveContext> {
  const res = await fetch(
    `${kernel.baseUrl}/context/issue?issueId=${issueId}`,
    {
      method: "PUT",
      headers: headers(kernel),
    }
  );

  if (!res.ok) {
    throw new Error("Failed to switch issue");
  }

  return res.json();
}

/**
 * DELETE /context/issue
 */
export async function clearIssue(
  kernel: KernelInfo
): Promise<ActiveContext> {
  const res = await fetch(`${kernel.baseUrl}/context/issue`, {
    method: "DELETE",
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to clear issue");
  }

  return res.json();
}

/**
 * DELETE /context
 */
export async function clearContext(
  kernel: KernelInfo
): Promise<void> {
  const res = await fetch(`${kernel.baseUrl}/context`, {
    method: "DELETE",
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to clear context");
  }
}
