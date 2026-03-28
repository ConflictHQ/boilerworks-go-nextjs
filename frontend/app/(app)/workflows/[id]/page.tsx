"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { api } from "@/lib/api";
import { toast } from "sonner";

interface WorkflowDefinition {
  uuid: string;
  name: string;
  description: string;
  status: string;
  states: { name: string; label: string; is_end?: boolean }[];
  transitions: { name: string; from: string; to: string }[];
}

export default function WorkflowDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const [workflow, setWorkflow] = useState<WorkflowDefinition | null>(null);

  useEffect(() => {
    api.get<WorkflowDefinition>(`/api/workflows/${id}`)
      .then(setWorkflow)
      .catch((e) => toast.error(e.message));
  }, [id]);

  if (!workflow) return <div className="p-6 text-sm text-muted-foreground">Loading...</div>;

  return (
    <div className="flex flex-1 flex-col gap-6 p-6">
      <div>
        <h1 className="text-xl font-semibold">{workflow.name}</h1>
        <p className="mt-1 text-sm text-muted-foreground">{workflow.description || "No description"}</p>
      </div>
      <div className="border-t border-border" />
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div className="rounded-xl border border-border bg-card p-6">
          <h3 className="mb-4 text-sm font-medium text-muted-foreground">States ({workflow.states?.length || 0})</h3>
          <div className="space-y-2">
            {workflow.states?.map((state) => (
              <div key={state.name} className="flex items-center justify-between rounded-lg border border-border px-4 py-2">
                <span className="font-medium">{state.label}</span>
                {state.is_end && <span className="text-xs text-muted-foreground">End State</span>}
              </div>
            ))}
          </div>
        </div>
        <div className="rounded-xl border border-border bg-card p-6">
          <h3 className="mb-4 text-sm font-medium text-muted-foreground">Transitions ({workflow.transitions?.length || 0})</h3>
          <div className="space-y-2">
            {workflow.transitions?.map((t) => (
              <div key={t.name} className="rounded-lg border border-border px-4 py-2">
                <span className="font-medium">{t.name}</span>
                <span className="ml-2 text-xs text-muted-foreground">{t.from} → {t.to}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
      <div className="flex items-center gap-2">
        <span className="inline-flex rounded-full bg-secondary px-2 py-0.5 text-xs capitalize">{workflow.status}</span>
      </div>
    </div>
  );
}
