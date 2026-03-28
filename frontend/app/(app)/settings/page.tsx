export default function SettingsPage() {
  return (
    <div className="flex flex-1 flex-col gap-6 p-6">
      <div>
        <h1 className="text-xl font-semibold">Settings</h1>
        <p className="mt-1 text-sm text-muted-foreground">Application settings.</p>
      </div>
      <div className="border-t border-border" />
      <div className="rounded-xl border border-border bg-card p-6">
        <p className="text-sm text-muted-foreground">Settings page is under development.</p>
      </div>
    </div>
  );
}
