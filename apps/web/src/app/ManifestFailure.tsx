import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";
import type { ManifestError } from "@/protocol/app-manifest";

export function ManifestFailure({ error }: { error: unknown }) {
  const manifestError = error as Partial<ManifestError>;
  const code = manifestError.code ?? "MANIFEST_LOAD_FAILED";
  const message =
    error instanceof Error ? error.message : "The application manifest could not be loaded.";
  return (
    <div className="min-h-screen bg-background px-6 py-12 text-foreground">
      <div className="mx-auto max-w-xl space-y-6">
        <div className="flex items-center justify-between gap-4">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground">
              schema-ui-core
            </p>
            <h1 className="mt-2 text-3xl font-semibold tracking-tight">Admin unavailable</h1>
          </div>
          <ThemeToggle />
        </div>
        <div className="border border-destructive/40 bg-card p-6">
          <p className="font-mono text-sm text-destructive">{code}</p>
          <p className="mt-3 text-sm leading-6 text-muted-foreground">{message}</p>
          <p className="mt-4 text-xs text-muted-foreground">
            The host refuses to render pages until a valid manifest is available.
          </p>
        </div>
        <Button type="button" variant="outline" onClick={() => window.location.reload()}>
          Retry manifest
        </Button>
      </div>
    </div>
  );
}
