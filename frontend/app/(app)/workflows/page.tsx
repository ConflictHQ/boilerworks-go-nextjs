"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import { Plus, Pencil, Trash2 } from "lucide-react";
import { toast } from "sonner";

interface WorkflowDefinition {
  uuid: string;
  name: string;
  description: string;
  status: string;
  states: { name: string; label: string }[];
  transitions: { name: string; from: string; to: string }[];
}

export default function WorkflowsPage() {
  const [workflows, setWorkflows] = useState<WorkflowDefinition[]>([]);
  const [loading, setLoading] = useState(true);

  const loadWorkflows = () => {
    api.get<{ workflows: WorkflowDefinition[] }>("/api/workflows")
      .then((data) => setWorkflows(data.workflows))
      .catch((e) => toast.error(e.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => { loadWorkflows(); }, []);

  const handleDelete = async (uuid: string) => {
    if (!confirm("Delete this workflow?")) return;
    try {
      await api.delete(`/api/workflows/${uuid}`);
      toast.success("Workflow deleted");
      loadWorkflows();
    } catch (e) {
      toast.error(e instanceof Error ? e.message : "Failed to delete");
    }
  };

  return (
    <div className="flex flex-1 flex-col gap-6 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold">Workflows</h1>
          <p className="mt-1 text-sm text-muted-foreground">Manage workflow definitions and instances.</p>
        </div>
        <Link href="/workflows/new" className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90">
          <Plus className="h-4 w-4" />New Workflow
        </Link>
      </div>
      <div className="border-t border-border" />
      {loading ? (
        <p className="text-sm text-muted-foreground">Loading...</p>
      ) : workflows.length === 0 ? (
        <p className="text-sm text-muted-foreground">No workflows yet.</p>
      ) : (
        <div className="overflow-hidden rounded-xl border border-border">
          <table className="w-full text-sm">
            <thead className="border-b border-border bg-muted/50">
              <tr>
                <th className="px-4 py-3 text-left font-medium">Name</th>
                <th className="px-4 py-3 text-left font-medium">Status</th>
                <th className="px-4 py-3 text-left font-medium">States</th>
                <th className="px-4 py-3 text-right font-medium">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {workflows.map((wf) => (
                <tr key={wf.uuid} className="hover:bg-muted/30">
                  <td className="px-4 py-3 font-medium">{wf.name}</td>
                  <td className="px-4 py-3"><span className="inline-flex rounded-full bg-secondary px-2 py-0.5 text-xs capitalize">{wf.status}</span></td>
                  <td className="px-4 py-3 text-muted-foreground">{wf.states?.length || 0} states</td>
                  <td className="px-4 py-3 text-right">
                    <div className="flex items-center justify-end gap-2">
                      <Link href={`/workflows/${wf.uuid}`} className="rounded p-1 text-muted-foreground hover:text-foreground"><Pencil className="h-4 w-4" /></Link>
                      <button onClick={() => handleDelete(wf.uuid)} className="rounded p-1 text-muted-foreground hover:text-destructive"><Trash2 className="h-4 w-4" /></button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
