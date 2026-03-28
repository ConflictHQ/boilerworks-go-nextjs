"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { api } from "@/lib/api";
import { toast } from "sonner";

interface Category {
  uuid: string;
  name: string;
  description: string;
}

export default function EditCategoryPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const [category, setCategory] = useState<Category | null>(null);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    api.get<Category>(`/api/categories/${id}`).then((c) => {
      setCategory(c);
      setName(c.name);
      setDescription(c.description);
    }).catch((e) => toast.error(e.message));
  }, [id]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      await api.put(`/api/categories/${id}`, { name, description });
      toast.success("Category updated");
      router.push("/categories");
    } catch (e) {
      toast.error(e instanceof Error ? e.message : "Failed to update");
    } finally {
      setSubmitting(false);
    }
  };

  if (!category) return <div className="p-6 text-sm text-muted-foreground">Loading...</div>;

  return (
    <div className="flex flex-1 flex-col gap-6 p-6">
      <div>
        <h1 className="text-xl font-semibold">Edit Category</h1>
        <p className="mt-1 text-sm text-muted-foreground">Update category details.</p>
      </div>
      <div className="border-t border-border" />
      <form onSubmit={handleSubmit} className="max-w-lg space-y-4">
        <div>
          <label className="mb-1 block text-sm font-medium">Name</label>
          <input type="text" value={name} onChange={(e) => setName(e.target.value)} required className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" />
        </div>
        <div>
          <label className="mb-1 block text-sm font-medium">Description</label>
          <textarea value={description} onChange={(e) => setDescription(e.target.value)} rows={3} className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" />
        </div>
        <div className="flex gap-3">
          <button type="submit" disabled={submitting} className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50">
            {submitting ? "Saving..." : "Save Changes"}
          </button>
          <button type="button" onClick={() => router.back()} className="rounded-lg border border-input px-4 py-2 text-sm font-medium hover:bg-accent">Cancel</button>
        </div>
      </form>
    </div>
  );
}
