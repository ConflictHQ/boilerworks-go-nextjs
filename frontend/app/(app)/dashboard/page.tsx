"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { Package, FolderTree, FileText, GitBranch } from "lucide-react";

interface DashboardData {
  items_by_status: Record<string, number>;
  category_count: number;
  form_count: number;
  submission_count: number;
  workflow_count: number;
  instance_count: number;
}

function StatCard({
  title,
  value,
  icon: Icon,
}: {
  title: string;
  value: number | string;
  icon: React.ComponentType<{ className?: string }>;
}) {
  return (
    <div className="rounded-xl border border-border bg-card p-6">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm text-muted-foreground">{title}</p>
          <p className="mt-1 text-2xl font-bold">{value}</p>
        </div>
        <Icon className="h-8 w-8 text-muted-foreground" />
      </div>
    </div>
  );
}

export default function DashboardPage() {
  const [data, setData] = useState<DashboardData | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api
      .get<DashboardData>("/api/dashboard")
      .then(setData)
      .catch((e) => setError(e.message));
  }, []);

  if (error) {
    return (
      <div className="flex flex-1 flex-col gap-6 p-6">
        <div>
          <h1 className="text-xl font-semibold">Dashboard</h1>
          <p className="mt-1 text-sm text-destructive">{error}</p>
        </div>
      </div>
    );
  }

  const totalItems = data
    ? Object.values(data.items_by_status).reduce((a, b) => a + b, 0)
    : 0;

  return (
    <div className="flex flex-1 flex-col gap-6 p-6">
      <div>
        <h1 className="text-xl font-semibold">Dashboard</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Overview of your application.
        </p>
      </div>
      <div className="border-t border-border" />
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard title="Items" value={totalItems} icon={Package} />
        <StatCard
          title="Categories"
          value={data?.category_count ?? 0}
          icon={FolderTree}
        />
        <StatCard
          title="Forms"
          value={data?.form_count ?? 0}
          icon={FileText}
        />
        <StatCard
          title="Workflows"
          value={data?.workflow_count ?? 0}
          icon={GitBranch}
        />
      </div>
      {data && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="rounded-xl border border-border bg-card p-6">
            <h3 className="mb-4 text-sm font-medium text-muted-foreground">
              Items by Status
            </h3>
            {Object.entries(data.items_by_status).length === 0 ? (
              <p className="text-sm text-muted-foreground">No items yet.</p>
            ) : (
              <div className="space-y-2">
                {Object.entries(data.items_by_status).map(
                  ([status, count]) => (
                    <div
                      key={status}
                      className="flex items-center justify-between"
                    >
                      <span className="text-sm capitalize">{status}</span>
                      <span className="text-sm font-medium">{count}</span>
                    </div>
                  ),
                )}
              </div>
            )}
          </div>
          <div className="rounded-xl border border-border bg-card p-6">
            <h3 className="mb-4 text-sm font-medium text-muted-foreground">
              Activity
            </h3>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm">Form Submissions</span>
                <span className="text-sm font-medium">
                  {data.submission_count}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm">Workflow Instances</span>
                <span className="text-sm font-medium">
                  {data.instance_count}
                </span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
