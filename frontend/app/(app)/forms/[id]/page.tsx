"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { api } from "@/lib/api";
import { toast } from "sonner";

interface FormDefinition {
  uuid: string;
  name: string;
  slug: string;
  description: string;
  status: string;
  schema: { name: string; label: string; type: string; required: boolean }[];
}

export default function FormDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const [form, setForm] = useState<FormDefinition | null>(null);

  useEffect(() => {
    api.get<FormDefinition>(`/api/forms/${id}`)
      .then(setForm)
      .catch((e) => toast.error(e.message));
  }, [id]);

  if (!form) return <div className="p-6 text-sm text-muted-foreground">Loading...</div>;

  return (
    <div className="flex flex-1 flex-col gap-6 p-6">
      <div>
        <h1 className="text-xl font-semibold">{form.name}</h1>
        <p className="mt-1 text-sm text-muted-foreground">{form.description || "No description"}</p>
      </div>
      <div className="border-t border-border" />
      <div className="rounded-xl border border-border bg-card p-6">
        <h3 className="mb-4 text-sm font-medium text-muted-foreground">Schema ({form.schema?.length || 0} fields)</h3>
        {form.schema && form.schema.length > 0 ? (
          <div className="space-y-2">
            {form.schema.map((field, i) => (
              <div key={i} className="flex items-center justify-between rounded-lg border border-border px-4 py-2">
                <div>
                  <span className="font-medium">{field.label}</span>
                  <span className="ml-2 text-xs text-muted-foreground">({field.type})</span>
                </div>
                {field.required && <span className="text-xs text-destructive">Required</span>}
              </div>
            ))}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground">No fields defined yet.</p>
        )}
      </div>
      <div className="flex items-center gap-2">
        <span className="inline-flex rounded-full bg-secondary px-2 py-0.5 text-xs capitalize">{form.status}</span>
        <span className="text-xs text-muted-foreground">Slug: {form.slug}</span>
      </div>
    </div>
  );
}
