import { useCallback, useEffect, useState } from "react";
import type { KernelInfo } from "../interfaces/kernel.interface.js";
import { getDailyPlan } from "../services/issue-api.js";
import type { DailyIssueSummary } from "../interfaces/issue.interface.js";

export function useDailyPlan(kernel: KernelInfo) {
  const [plan, setPlan] = useState<DailyIssueSummary | null>(null);

  const reload = useCallback(async () => {
    try {
      const data = await getDailyPlan(kernel);
      setPlan(data);
    } catch (err) {
      console.error("Failed to load daily plan", err);
    }
  }, [kernel.baseUrl]);

  // initial load
  useEffect(() => {
    reload();
  }, [reload]);

  return { plan, reload };
}
