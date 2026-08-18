// Import CSV template download (W16-F03): a visible entry point on the users
// page so operators can see the exact header columns before uploading.
import { useState } from "react";

import { useAuth } from "@/account/AuthContext";
import { Button } from "@/components/ui/button";
import { useTranslate } from "@/i18n/runtime";
import {
  registerCustomComponent,
  type CustomComponentProps,
} from "@/renderer/custom-components";

export function ImportTemplateDownload(_props: CustomComponentProps) {
  const t = useTranslate();
  const { authFetch } = useAuth();
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const download = async () => {
    if (busy) {
      return;
    }
    setBusy(true);
    setError(null);
    try {
      const response = await authFetch("/api/import/users/template", {
        headers: { Accept: "text/csv" },
      });
      if (!response.ok) {
        setError(t("importTemplate.failed"));
        return;
      }
      const blob = await response.blob();
      const objectUrl = URL.createObjectURL(blob);
      const anchor = document.createElement("a");
      anchor.href = objectUrl;
      anchor.download = "users-import-template.csv";
      document.body.appendChild(anchor);
      anchor.click();
      anchor.remove();
      URL.revokeObjectURL(objectUrl);
    } catch {
      setError(t("importTemplate.failed"));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="space-y-1">
      <Button
        type="button"
        variant="outline"
        disabled={busy}
        className="h-8 text-xs"
        onClick={() => void download()}
        data-import-template-download
      >
        {t("importTemplate.download")}
      </Button>
      {error !== null ? (
        <p role="alert" className="text-xs text-destructive" data-import-template-error>
          {error}
        </p>
      ) : null}
    </div>
  );
}

registerCustomComponent("import-template-download", ImportTemplateDownload);
