import type { KernelInfo } from "../interfaces/kernel.interface.js";
import type { Issue } from "../interfaces/issue.interface.js";

function headers(kernel: KernelInfo) {
  return {
    Authorization: `Bearer ${kernel.token}`,
    "Content-Type": "application/json",
  };
}

/**
 * GET /issues?streamId=...
 */
export async function getIssues(
  kernel: KernelInfo,
  streamId?: string
): Promise<Issue[]> {
  if (!streamId) {
    throw new Error("streamId required");
  }

  const res = await fetch(
    `${kernel.baseUrl}/issues?streamId=${streamId}`,
    {
      headers: headers(kernel),
    }
  );

  if (!res.ok) {
    throw new Error("Failed to fetch issues");
  }

  return res.json();
}

/**
 * GET /issues/summary/today
 */
export async function getDailyPlan(kernel: KernelInfo) {
  const res = await fetch(
    `${kernel.baseUrl}/issues/summary/today`,
    {
      headers: headers(kernel),
    }
  );

  if (!res.ok) {
    throw new Error("Failed to fetch daily plan");
  }

  return res.json();
}

/**
 * POST /issue
 */
export async function createIssueApi(
  kernel: KernelInfo,
  streamId: string,
  title: string,
  estimateMinutes?: number
) {
  const res = await fetch(`${kernel.baseUrl}/issue`, {
    method: "POST",
    headers: headers(kernel),
    body: JSON.stringify({
      streamId,
      title,
      estimateMinutes,
    }),
  });

  if (!res.ok) {
    throw new Error("Failed to create issue");
  }

  return res.json();
}

/**
 * PUT /issue/:id
 */
export async function updateIssueApi(
  kernel: KernelInfo,
  id: string,
  data: {
    title?: string;
    estimateMinutes?: number;
    notes?: string;
  }
) {
  const res = await fetch(`${kernel.baseUrl}/issue/${id}`, {
    method: "PUT",
    headers: headers(kernel),
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    throw new Error("Failed to update issue");
  }

  return res.json();
}

/**
 * DELETE /issue/:id
 */
export async function deleteIssueApi(
  kernel: KernelInfo,
  id: string
) {
  const res = await fetch(`${kernel.baseUrl}/issue/${id}`, {
    method: "DELETE",
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to delete issue");
  }
}

/**
 * PUT /issue/:id/status
 */
export async function changeIssueStatusApi(
  kernel: KernelInfo,
  id: string,
  status: string
) {
  const res = await fetch(`${kernel.baseUrl}/issue/${id}/status`, {
    method: "PUT",
    headers: headers(kernel),
    body: JSON.stringify({ status }),
  });

  if (!res.ok) {
    throw new Error("Failed to change issue status");
  }

  return res.json();
}

/**
 * PUT /issue/:id/todo
 */
export async function markIssueTodoToday(
  kernel: KernelInfo,
  id: string
) {
  const res = await fetch(`${kernel.baseUrl}/issue/${id}/todo`, {
    method: "PUT",
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to mark todo");
  }

  return res.json();
}

/**
 * PUT /issue/:id/todo/clear
 */
export async function clearIssueTodoToday(
  kernel: KernelInfo,
  id: string
) {
  const res = await fetch(`${kernel.baseUrl}/issue/${id}/todo/clear`, {
    method: "PUT",
    headers: headers(kernel),
  });

  if (!res.ok) {
    throw new Error("Failed to clear todo");
  }

  return res.json();
}
