"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import { Plus, Pencil, Trash2 } from "lucide-react";
import { toast } from "sonner";

interface FormDefinition {
  uuid: string;
  name: string;
  slug: string;
  description: string;
  status: string;
  created_at: string;
}

export default function FormsPage() {
  const [forms, setForms] = useState<FormDefinition[]>([]);
  const [loading, setLoading] = useState(true);

  const loadForms = () => {
    api.get<{ forms: FormDefinition[] }>("/api/forms")
      .then((data) => setForms(data.forms))
      .catch((e) => toast.error(e.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => { loadForms(); }, []);

  const handleDelete = async (uuid: string) => {
    if (!confirm("Delete this form definition?")) return;
    try {
      await api.delete(`/api/forms/${uuid}`);
      toast.success("Form deleted");
      loadForms();
    } catch (e) {
      toast.error(e instanceof Error ? e.message : "Failed to delete");
    }
  };

  return (
    <div className="flex flex-1 flex-col gap-6 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold">Forms</h1>
          <p className="mt-1 text-sm text-muted-foreground">Manage form definitions and submissions.</p>
        </div>
        <Link href="/forms/new" className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90">
          <Plus className="h-4 w-4" />New Form
        </Link>
      </div>
      <div className="border-t border-border" />
      {loading ? (
        <p className="text-sm text-muted-foreground">Loading...</p>
      ) : forms.length === 0 ? (
        <p className="text-sm text-muted-foreground">No forms yet.</p>
      ) : (
        <div className="overflow-hidden rounded-xl border border-border">
          <table className="w-full text-sm">
            <thead className="border-b border-border bg-muted/50">
              <tr>
                <th className="px-4 py-3 text-left font-medium">Name</th>
                <th className="px-4 py-3 text-left font-medium">Slug</th>
                <th className="px-4 py-3 text-left font-medium">Status</th>
                <th className="px-4 py-3 text-right font-medium">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {forms.map((form) => (
                <tr key={form.uuid} className="hover:bg-muted/30">
                  <td className="px-4 py-3 font-medium">{form.name}</td>
                  <td className="px-4 py-3 text-muted-foreground">{form.slug}</td>
                  <td className="px-4 py-3"><span className="inline-flex rounded-full bg-secondary px-2 py-0.5 text-xs capitalize">{form.status}</span></td>
                  <td className="px-4 py-3 text-right">
                    <div className="flex items-center justify-end gap-2">
                      <Link href={`/forms/${form.uuid}`} className="rounded p-1 text-muted-foreground hover:text-foreground"><Pencil className="h-4 w-4" /></Link>
                      <button onClick={() => handleDelete(form.uuid)} className="rounded p-1 text-muted-foreground hover:text-destructive"><Trash2 className="h-4 w-4" /></button>
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
