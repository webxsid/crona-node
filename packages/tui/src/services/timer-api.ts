import type { KernelInfo } from "../interfaces/kernel.interface.js";

function headers(kernel: KernelInfo) {
  return {
    Authorization: `Bearer ${kernel.token}`,
    "Content-Type": "application/json",
  };
}

/**
 * GET /timer/state
 */
export async function getTimerState(kernel: KernelInfo) {
  const res = await fetch(`${kernel.baseUrl}/timer/state`, {
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to fetch timer state");
  }

  return res.json();
}

/**
 * POST /timer/start?issueId=...
 */
export async function startTimer(
  kernel: KernelInfo,
  issueId?: string
) {
  const url = issueId
    ? `${kernel.baseUrl}/timer/start?issueId=${issueId}`
    : `${kernel.baseUrl}/timer/start`;

  const res = await fetch(url, {
    method: "POST",
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to start timer");
  }

  return res.json();
}

/**
 * POST /timer/pause
 */
export async function pauseTimer(kernel: KernelInfo) {
  const res = await fetch(`${kernel.baseUrl}/timer/pause`, {
    method: "POST",
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to pause timer");
  }

  return res.json();
}

/**
 * POST /timer/resume
 */
export async function resumeTimer(kernel: KernelInfo) {
  const res = await fetch(`${kernel.baseUrl}/timer/resume`, {
    method: "POST",
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to resume timer");
  }

  return res.json();
}

/**
 * POST /timer/end
 */
export async function endTimer(
  kernel: KernelInfo,
  commitMessage?: string
) {
  const res = await fetch(`${kernel.baseUrl}/timer/end`, {
    method: "POST",
    headers: headers(kernel),
    body: JSON.stringify({ commitMessage }),
  });

  if (!res.ok) {
    throw new Error("Failed to end timer");
  }

  return res.json();
}
