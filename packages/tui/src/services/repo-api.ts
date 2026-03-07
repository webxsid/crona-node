import type { KernelInfo } from "../interfaces/kernel.interface.js";
import type { Repo } from "../interfaces/repo.interface.js";

function headers(kernel: KernelInfo) {
  return {
    Authorization: `Bearer ${kernel.token}`,
    "Content-Type": "application/json",
  };
}

/**
 * GET /repos
 */
export async function getRepos(kernel: KernelInfo): Promise<Repo[]> {
  const res = await fetch(`${kernel.baseUrl}/repos`, {
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error(`Failed to fetch repos (${res.status})`);
  }

  return res.json();
}

/**
 * POST /commands/repo
 */
export async function createRepoApi(
  kernel: KernelInfo,
  name: string
): Promise<Repo> {
  const res = await fetch(`${kernel.baseUrl}/commands/repo`, {
    method: "POST",
    headers: headers(kernel),
    body: JSON.stringify({ name }),
  });

  if (!res.ok) {
    throw new Error("Failed to create repo");
  }

  return res.json();
}

/**
 * PUT /commands/repo/:id
 */
export async function updateRepoApi(
  kernel: KernelInfo,
  id: string,
  name: string
): Promise<Repo> {
  const res = await fetch(`${kernel.baseUrl}/commands/repo/${id}`, {
    method: "PUT",
    headers: headers(kernel),
    body: JSON.stringify({ name }),
  });

  if (!res.ok) {
    throw new Error("Failed to update repo");
  }

  return res.json();
}

/**
 * DELETE /commands/repo/:id
 */
export async function deleteRepoApi(
  kernel: KernelInfo,
  id: string
): Promise<void> {
  const res = await fetch(`${kernel.baseUrl}/commands/repo/${id}`, {
    method: "DELETE",
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to delete repo");
  }
}
